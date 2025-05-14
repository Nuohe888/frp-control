package cron

import (
	"easy-fiber-admin/model/frp"
	"easy-fiber-admin/module/frp/internal/pkg/client"
	"easy-fiber-admin/module/frp/internal/service"
	"easy-fiber-admin/pkg/logger"
	"easy-fiber-admin/pkg/sql"
	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

func Init() {
	cronScheduler = cron.New()
	// 添加每分钟运行一次的flowCheckTask
	_, err := cronScheduler.AddFunc("*/1 * * * *", flowCheckTask)
	if err != nil {
		logger.Get().Errorf("Failed to schedule flowCheckTask: %v", err)
		return
	}

	// 启动定时任务
	cronScheduler.Start()
}

func flowCheckTask() {
	logger.Get().Info("流量检查任务开始执行...")

	var node []frp.Node
	sql.Get().Find(&node)
	var list []client.FlowRes
	for _, v := range node {
		if v.Ip == "" || *v.Status != 1 {
			continue
		}
		flow, err := client.Flow(v.Ip)
		if err != nil {
			logger.Get().Error(err)
			continue
		}
		list = append(list, *flow)
	}

	service.GetUserSrv().DeductFlow(list)
	runIdList := service.GetUserSrv().GetFlowShortage()
	for _, v := range runIdList {
		go func() {
			err := client.Kill(v.Ip, v.RunId)
			if err != nil {
				logger.Get().Error(err)
			}
		}()
	}
	logger.Get().Info("流量检查任务结束执行...")
}
