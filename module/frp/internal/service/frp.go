package service

import (
	"easy-fiber-admin/model/frp"
	"easy-fiber-admin/module/frp/internal/vo"
	"easy-fiber-admin/pkg/logger"
	"easy-fiber-admin/pkg/sql"
	"errors"
	"gorm.io/gorm"
)

type FrpSrv struct {
	db  *gorm.DB
	log logger.ILog
}

var frpSrv *FrpSrv

func InitFrpSrv() {
	frpSrv = &FrpSrv{
		db:  sql.Get(),
		log: logger.Get(),
	}
}

func GetFrpSrv() *FrpSrv {
	if frpSrv == nil {
		panic("service frp init failed")
	}
	return frpSrv
}

func (i *FrpSrv) Auth(req *vo.AuthReq) error {
	var order frp.Order
	i.db.Where("node_name=?", req.Node).
		Where("user_uuid=?", req.Uuid).
		Find(&order)
	if order.Id == 0 {
		return errors.New("不存在该用户")
	}
	var user frp.User
	i.db.Where("uuid=?", req.Uuid).Find(&user)
	if user.Uuid == "" || user.Status != 1 {
		return errors.New("用户被封禁")
	}
	if user.Flow < 0 {
		return errors.New("流量已超出")
	}
	order.RunId = req.RunId
	err := i.db.Save(&order).Error
	if err != nil {
		i.log.Debug(err)
		return errors.New("服务端出错")
	}
	return nil
}

func (i *FrpSrv) Auth2(req *vo.Auth2Req) error {
	var order frp.Order
	i.db.Where("node_name=?", req.Node).
		Where("username=?", req.Name).
		Where("port=?", req.Port).
		Where("type=?", req.Type).
		Where("user_uuid=?", req.Uuid).Find(&order)
	if order.Id == 0 {
		return errors.New("没有开通隧道,或隧道信息填写错误")
	}
	return nil
}

func (i *FrpSrv) Speed(req *vo.SpeedReq) *vo.SpeedRes {
	var user frp.User
	i.db.Where("uuid=?", req.Uuid).
		Find(&user)

	return &vo.SpeedRes{
		Limit: user.Speed,
	}
}
