package service

import (
	"easy-fiber-admin/model/frp"
	"easy-fiber-admin/module/frp/internal/utils"
	"easy-fiber-admin/module/frp/internal/vo"
	"easy-fiber-admin/pkg/logger"
	"easy-fiber-admin/pkg/sql"
	"errors"
	"gorm.io/gorm"
)

type OrderSrv struct {
	db  *gorm.DB
	log logger.ILog
}

var orderSrv *OrderSrv

func InitOrderSrv() {
	orderSrv = &OrderSrv{
		db:  sql.Get(),
		log: logger.Get(),
	}
}

func GetOrderSrv() *OrderSrv {
	if orderSrv == nil {
		panic("service order init failed")
	}
	return orderSrv
}

func (i *OrderSrv) Add(order *frp.Order) error {
	node := GetNodeSrv().Get(order.NodeId)
	if node.Id == 0 {
		return errors.New("节点id出错")
	}
	user := GetUserSrv().Get(order.UserUuid)
	if user.Uuid == "" {
		return errors.New("用户uuid出错")
	}
	order.NodeName = node.Name
	order.NodeIp = node.Ip
	if order.Speed == 0 {
		order.Speed = user.Speed
	}
	return i.db.Create(&order).Error
}

func (i *OrderSrv) Del(id string) error {
	return i.db.Where("id = ?", id).Delete(&frp.Order{}).Error
}

func (i *OrderSrv) Put(id string, order *frp.Order) error {
	var _order frp.Order
	i.db.Where("id = ?", id).Find(&_order)
	if _order.Id == 0 {
		return errors.New("不存在该Id")
	}
	//originalStatus := order.Status
	utils.MergeStructs(&_order, order)
	//_order.Status = originalStatus
	return i.db.Save(&_order).Error
}

func (i *OrderSrv) Get(id string) frp.Order {
	var order frp.Order
	i.db.Where("id = ?", id).Find(&order)
	return order
}

func (i *OrderSrv) List(page, limit int) *vo.List {
	var items []frp.Order
	var total int64
	if limit == 0 {
		limit = 20
	}
	db := i.db
	i.db.Limit(limit).Offset((page - 1) * limit).Find(&items)
	db.Model(&frp.Order{}).Count(&total)
	return &vo.List{
		Items: items,
		Total: total,
	}
}
