package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var app *fiber.App
var port int

func Init(_port int) {
	app = newFiber()
	port = _port

	// 添加静态文件服务
	setupStaticFiles()
}

// 设置静态文件服务
func setupStaticFiles() {
	// 首先尝试使用嵌入的静态文件
	if setupEmbeddedStaticFiles() {
		return
	}

	// 如果嵌入的静态文件不可用，尝试从本地文件系统提供
	if _, err := os.Stat("dist"); !os.IsNotExist(err) {
		// 使用filesystem中间件提供静态文件服务
		app.Use("/", filesystem.New(filesystem.Config{
			Root:       http.Dir("dist"),
			Browse:     false,
			Index:      "index.html",
			MaxAge:     3600,
			PathPrefix: "",
		}))

		// 处理SPA路由，将所有未匹配的路由重定向到index.html
		app.Use(func(c *fiber.Ctx) error {
			path := c.Path()

			// 如果请求的是API路径，跳过处理
			if len(path) >= 4 && path[:4] == "/api" {
				return c.Next()
			}

			// 检查请求的文件是否存在
			if _, err := os.Stat(filepath.Join("dist", path)); os.IsNotExist(err) {
				// 文件不存在，返回index.html
				return c.SendFile("dist/index.html")
			}

			return c.Next()
		})

		log.Println("静态文件服务已启用，前端文件将从本地 'dist' 目录提供")
	} else {
		log.Println("未找到 'dist' 目录，静态文件服务未启用")
	}
}

func Get() *fiber.App {
	return app
}

func Start() {
	app.Listen(fmt.Sprintf(":%d", port))
}

func Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭服务器中 ...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("服务器退出失败:", err)
	}
	log.Println("服务器已关闭...")
}
