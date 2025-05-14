package service

import (
	"easy-fiber-admin/model/frp"
	"easy-fiber-admin/module/frp/internal/pkg/client"
	"easy-fiber-admin/module/frp/internal/utils"
	"easy-fiber-admin/module/frp/internal/vo"
	"easy-fiber-admin/pkg/logger"
	"easy-fiber-admin/pkg/sql"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSrv struct {
	db  *gorm.DB
	log logger.ILog
}

var userSrv *UserSrv

func InitUserSrv() {
	userSrv = &UserSrv{
		db:  sql.Get(),
		log: logger.Get(),
	}
}

func GetUserSrv() *UserSrv {
	if userSrv == nil {
		panic("service user init failed")
	}
	return userSrv
}

func (i *UserSrv) Add(user *frp.User) error {
	user.Uuid = uuid.NewString()
	user.Flow *= 1024 * 1024
	return i.db.Create(&user).Error
}

func (i *UserSrv) Del(uuid string) error {
	return i.db.Where("uuid = ?", uuid).Delete(&frp.User{}).Error
}

func (i *UserSrv) Put(uuid string, user *frp.User) error {
	var _user frp.User
	result := i.db.Where("uuid = ?", uuid).Find(&_user)
	if result.RowsAffected == 0 {
		return errors.New("不存在该Uuid")
	}
	originalStatus := user.Status
	if _user.Uuid != uuid {
		return errors.New("不允许修改UUID")
	}
	_user.Status = originalStatus
	user.Flow = _user.Flow + user.Flow*1024*1024

	utils.MergeStructs(&_user, user)

	return i.db.Save(&_user).Error
}

func (i *UserSrv) Get(uuid string) frp.User {
	var user frp.User
	i.db.Where("uuid = ?", uuid).Find(&user)
	return user
}

func (i *UserSrv) List(page, limit int) *vo.List {
	var items []frp.User
	var total int64
	if limit == 0 {
		limit = 20
	}
	db := i.db
	i.db.Limit(limit).Offset((page - 1) * limit).Find(&items)
	db.Model(&frp.User{}).Count(&total)
	return &vo.List{
		Items: items,
		Total: total,
	}
}

// 存储上一次的流量统计数据
var previousFlowStats = make(map[string]int64)

// DeductFlow 扣除用户流量，接收流量使用数据列表
func (i *UserSrv) DeductFlow(flowList []client.FlowRes) {
	// 创建一个map来汇总当前的流量统计
	currentFlowStats := make(map[string]int64)
	// 创建一个map来存储用户名的增量流量
	usernameIncrementalFlow := make(map[string]int64)

	// 汇总当前的流量统计
	for _, flowRes := range flowList {
		// 处理各种类型的代理
		processProxies := func(proxies []*client.ProxyStatsInfo) {
			for _, proxy := range proxies {
				if proxy.Name == "" {
					continue
				}

				// 当前总流量
				totalFlow := proxy.TodayTrafficIn + proxy.TodayTrafficOut

				// 计算增量流量（无论代理是否在线）
				previousFlow := previousFlowStats[proxy.Name]
				incrementalFlow := totalFlow - previousFlow

				// 只有当增量为正时才计入
				if incrementalFlow > 0 {
					username := proxy.Name
					usernameIncrementalFlow[username] += incrementalFlow
				}

				// 如果代理在线，更新当前统计；否则，从之前的统计中移除
				if proxy.Status == "online" {
					currentFlowStats[proxy.Name] = totalFlow
				} else {
					delete(previousFlowStats, proxy.Name)
				}
			}
		}

		// 处理各种类型的代理
		processProxies(flowRes.Tcp)
		processProxies(flowRes.Udp)
		processProxies(flowRes.Https)
		processProxies(flowRes.Tcpmux)
		processProxies(flowRes.Stcp)
		processProxies(flowRes.Sudp)
	}

	// 更新之前的流量统计
	for name, flow := range currentFlowStats {
		previousFlowStats[name] = flow
	}

	// 处理每个用户名的流量扣除
	for username, usedFlow := range usernameIncrementalFlow {
		if usedFlow <= 0 {
			continue
		}

		// 开始数据库事务
		tx := i.db.Begin()

		// 查询订单信息
		var order frp.Order
		if err := tx.Where("username = ?", username).Find(&order).Error; err != nil {
			tx.Rollback()
			i.log.Errorf("查询订单失败 username: %s, error: %v", username, err)
			continue
		}

		// 检查订单是否存在
		if order.Id == 0 {
			tx.Rollback()
			i.log.Warnf("未找到订单 username: %s", username)
			continue
		}

		// 查询用户信息
		var user frp.User
		if err := tx.Where("uuid = ?", order.UserUuid).Find(&user).Error; err != nil {
			tx.Rollback()
			i.log.Errorf("查询用户失败 uuid: %s, error: %v", order.UserUuid, err)
			continue
		}

		// 检查用户是否存在
		if user.Uuid == "" {
			tx.Rollback()
			i.log.Warnf("未找到用户 uuid: %s", order.UserUuid)
			continue
		}

		// 扣除用户流量
		user.Flow -= usedFlow

		// 更新用户信息
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			i.log.Errorf("更新用户流量失败 uuid: %s, error: %v", user.Uuid, err)
			continue
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			i.log.Errorf("提交事务失败: %v", err)
			continue
		}

	}
}

type GetFlowShortage struct {
}

// GetFlowShortageRes 定义流量不足用户的隧道信息
type GetFlowShortageRes struct {
	RunId string `json:"runId"`
	Ip    string `json:"ip"`
}

// GetFlowShortage 获取流量不足的用户的在线隧道信息
func (i *UserSrv) GetFlowShortage() []GetFlowShortageRes {
	var result []GetFlowShortageRes

	// 查询流量不足的用户
	var users []frp.User
	if err := i.db.Where("flow < 0").Find(&users).Error; err != nil {
		i.log.Errorf("查询流量不足用户失败: %v", err)
		return result
	}

	// 如果没有流量不足的用户，直接返回空数组
	if len(users) == 0 {
		return result
	}

	// 收集所有流量不足用户的UUID
	var userUuids []string
	for _, user := range users {
		userUuids = append(userUuids, user.Uuid)
	}

	// 查询这些用户的所有订单
	var orders []frp.Order
	if err := i.db.Where("user_uuid IN ?", userUuids).Find(&orders).Error; err != nil {
		i.log.Errorf("查询用户订单失败: %v", err)
		return result
	}

	// 收集所有有效的隧道信息
	for _, order := range orders {
		// 只收集有runId的订单（表示隧道在线）
		if order.RunId != "" {
			result = append(result, GetFlowShortageRes{
				RunId: order.RunId,
				Ip:    order.NodeIp,
			})
			i.log.Infof("用户 %s 流量不足，将关闭隧道 runId: %s, nodeIp: %s",
				order.UserUuid, order.RunId, order.NodeIp)
		}
	}

	return result
}
