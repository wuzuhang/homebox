package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"home-box/config"
	"home-box/cron"
	"home-box/pkg/storage"
	"home-box/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置文件
	config.InitConfig()

	// 2. 设置 Gin 运行模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 3. 初始化数据库
	config.InitDB()
	// 初始化储存桶
	bucketName := config.AppConfig.COS.Bucket
	region := config.AppConfig.COS.Region
	secretID := config.AppConfig.COS.SecretID
	secretKey := config.AppConfig.COS.SecretKey
	storage.InitCOS(bucketName, region, secretID, secretKey)

	// fileService := service.NewFileService(cosClient)

	// 4. 初始化定时任务调度器 (必须在启动阻塞式 Web 服务之前初始化并 Start)
	cronMgr := cron.GetManager()
	_, err := cronMgr.AddJob("0 0 0 * * *", "每日更新药品库存", func() {
		cron.UpdateMedicineStock()
	})
	if err != nil {
		log.Printf("添加定时任务失败: %v\n", err)
	}
	cronMgr.Start() // 启动 Cron（后台 goroutine 运行）

	// 5. 初始化路由
	r := router.InitRouter()
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 6. 在单独的 Goroutine 中启动 HTTP 服务（防止阻塞主线程）
	go func() {
		fmt.Printf("🚀 服务已启动，监听地址: http://localhost%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务器异常退出: %v\n", err)
		}
	}()
	// 7. 监听退出信号（优雅停机）
	quit := make(chan os.Signal, 1)
	// 捕获 SIGINT (Ctrl+C) 和 SIGTERM (kill 信号)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务...")

	// 8. 关闭 Cron 调度器（等待正在运行的任务执行完毕）
	cronMgr.Stop()

	// 9. 优雅关闭 Gin HTTP 服务（给未完成的 HTTP 请求 5 秒缓冲时间）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP 服务器强制关闭: %v\n", err)
	}

	log.Println("服务已完全安全退出")
}
