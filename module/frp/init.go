package frp

import (
	"easy-fiber-admin/module/frp/internal/controller"
	"easy-fiber-admin/module/frp/internal/cron"
	"easy-fiber-admin/module/frp/internal/service"
)

func Init() {
	service.Init()
	controller.Init()
	cron.Init()
}
