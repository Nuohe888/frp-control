package controller

import (
	"easy-fiber-admin/model/frp"
	"easy-fiber-admin/module/frp/internal/service"
	"easy-fiber-admin/module/frp/internal/vo"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type orderCtl struct {
	srv *service.OrderSrv
}

var OrderCtl *orderCtl

func InitOrderCtl() {
	OrderCtl = &orderCtl{
		srv: service.GetOrderSrv(),
	}
}

func (i *orderCtl) Add(c *fiber.Ctx) error {
	var req frp.Order
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Add(&req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) Del(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := i.srv.Del(id); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	var req frp.Order
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Put(id, &req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) Get(c *fiber.Ctx) error {
	id := c.Query("id")
	return vo.ResultOK(i.srv.Get(id), c)
}

func (i *orderCtl) List(c *fiber.Ctx) error {
	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	return vo.ResultOK(i.srv.List(pageInt, limitInt), c)
}
