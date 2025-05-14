package controller

import (
	"easy-fiber-admin/module/frp/internal/service"
	"easy-fiber-admin/module/frp/internal/vo"
	"github.com/gofiber/fiber/v2"
)

type frpCtl struct {
	srv *service.FrpSrv
}

var FrpCtl *frpCtl

func InitFrpCtl() {
	FrpCtl = &frpCtl{
		srv: service.GetFrpSrv(),
	}
}

func (i *frpCtl) Auth(c *fiber.Ctx) error {
	var req vo.AuthReq
	err := vo.BodyParser(&req, c)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	err = i.srv.Auth(&req)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *frpCtl) Auth2(c *fiber.Ctx) error {
	var req vo.Auth2Req
	err := vo.BodyParser(&req, c)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	err = i.srv.Auth2(&req)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *frpCtl) Speed(c *fiber.Ctx) error {
	var req vo.SpeedReq
	err := vo.BodyParser(&req, c)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(i.srv.Speed(&req), c)
}
