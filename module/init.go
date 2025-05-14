package module

import (
	"easy-fiber-admin/module/frp"
	"easy-fiber-admin/module/system"
)

func Init() {
	system.Init()
	frp.Init()
}
