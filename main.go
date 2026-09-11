package main

import (
	"fmt"
	"home-box/config"
	"home-box/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置文件
	config.InitConfig()

	// 2. 根据配置设置 Gin 运行模式 (debug / release)
	gin.SetMode(config.AppConfig.Server.Mode)

	// 3. 初始化数据库
	config.InitDB()

	// 4. 初始化路由
	r := router.InitRouter()

	// 5. 从配置中读取端口并启动服务
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	fmt.Printf("🚀 服务已启动，监听地址: http://localhost%s\n", addr)

	if err := r.Run(addr); err != nil {
		panic("服务器启动失败: " + err.Error())
	}
}
