package server

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed all:dist
var distFS embed.FS

// setupEmbeddedStaticFiles 设置嵌入的静态文件服务
func setupEmbeddedStaticFiles() bool {
	// 尝试从嵌入的文件系统提供静态文件
	distSubFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Printf("无法加载嵌入的dist目录: %v", err)
		return false
	}

	// 使用filesystem中间件提供嵌入的静态文件服务
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(distSubFS),
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

		// 返回index.html以支持SPA路由
		return c.SendFile("index.html")
	})

	log.Println("静态文件服务已启用，前端文件将从嵌入的 'dist' 目录提供")
	return true
}
