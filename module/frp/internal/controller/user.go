package controller

import (
	"easy-fiber-admin/model/frp"
	"easy-fiber-admin/module/frp/internal/service"
	"easy-fiber-admin/module/frp/internal/vo"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type userCtl struct {
	srv *service.UserSrv
}

var UserCtl *userCtl

func InitUserCtl() {
	UserCtl = &userCtl{
		srv: service.GetUserSrv(),
	}
}

func (i *userCtl) Add(c *fiber.Ctx) error {
	var req frp.User
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Add(&req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *userCtl) Del(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if err := i.srv.Del(uuid); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *userCtl) Put(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	var req frp.User
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Put(uuid, &req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *userCtl) Get(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	return vo.ResultOK(i.srv.Get(uuid), c)
}

func (i *userCtl) List(c *fiber.Ctx) error {
	page := c.Query("page")
	limit := c.Query("limit")
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	return vo.ResultOK(i.srv.List(pageInt, limitInt), c)
}
