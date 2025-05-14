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

type NodeSrv struct {
	db  *gorm.DB
	log logger.ILog
}

var nodeSrv *NodeSrv

func InitNodeSrv() {
	nodeSrv = &NodeSrv{
		db:  sql.Get(),
		log: logger.Get(),
	}
}

func GetNodeSrv() *NodeSrv {
	if nodeSrv == nil {
		panic("service node init failed")
	}
	return nodeSrv
}

func (i *NodeSrv) Add(node *frp.Node) error {
	return i.db.Create(&node).Error
}

func (i *NodeSrv) Del(id string) error {
	return i.db.Where("id = ?", id).Delete(&frp.Node{}).Error
}

func (i *NodeSrv) Put(id string, node *frp.Node) error {
	var _node frp.Node
	i.db.Where("id = ?", id).Find(&_node)
	if _node.Id == 0 {
		return errors.New("不存在该Id")
	}
	originalStatus := node.Status
	utils.MergeStructs(&_node, node)
	_node.Status = originalStatus
	return i.db.Save(&_node).Error
}

func (i *NodeSrv) Get(id any) frp.Node {
	var node frp.Node
	i.db.Where("id = ?", id).Find(&node)
	return node
}

func (i *NodeSrv) List(page, limit int) *vo.List {
	var items []frp.Node
	var total int64
	if limit == 0 {
		limit = 20
	}
	db := i.db
	i.db.Limit(limit).Offset((page - 1) * limit).Find(&items)
	db.Model(&frp.Node{}).Count(&total)
	return &vo.List{
		Items: items,
		Total: total,
	}
}

func (i *NodeSrv) ListAll() []frp.Node {
	var list []frp.Node
	i.db.Find(&list)
	return list
}
