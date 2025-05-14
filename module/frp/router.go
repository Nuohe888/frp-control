package frp

import (
	"easy-fiber-admin/module/frp/internal/controller"
	"easy-fiber-admin/module/frp/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(r fiber.Router) {
	r.Post("/node/auth", controller.FrpCtl.Auth)
	r.Post("/node/auth2", controller.FrpCtl.Auth2)
	r.Post("/node/speed", controller.FrpCtl.Speed)

	auth := r.Group("/admin/frp")
	auth.Use(middleware.JWT()).
		Use(middleware.Casbin())

	auth.Post("node", controller.NodeCtl.Add)
	auth.Delete("node/:id", controller.NodeCtl.Del)
	auth.Put("node/:id", controller.NodeCtl.Put)
	auth.Get("node", controller.NodeCtl.Get)
	auth.Get("node/list", controller.NodeCtl.List)
	auth.Get("node/list/all", controller.NodeCtl.ListAll)

	auth.Post("order", controller.OrderCtl.Add)
	auth.Delete("order/:id", controller.OrderCtl.Del)
	auth.Put("order/:id", controller.OrderCtl.Put)
	auth.Get("order", controller.OrderCtl.Get)
	auth.Get("order/list", controller.OrderCtl.List)

	auth.Post("user", controller.UserCtl.Add)
	auth.Delete("user/:uuid", controller.UserCtl.Del)
	auth.Put("user/:uuid", controller.UserCtl.Put)
	auth.Get("user", controller.UserCtl.Get)
	auth.Get("user/list", controller.UserCtl.List)
}
