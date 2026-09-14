package cron

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron *cron.Cron
	once sync.Once
}

var instance *CronService

// GetManager 获取定时任务管理器单例
func GetManager() *CronService {
	if instance == nil {
		instance = &CronService{}
		instance.init()
	}
	return instance
}

// 初始化 Cron 实例
func (s *CronService) init() {
	s.once.Do(func() {
		// 1. 设置自定义日志输出格式
		logger := cron.VerbosePrintfLogger(log.New(log.Writer(), "[CRON] ", log.LstdFlags))

		// 2. 创建支持秒级表达式的 Cron 实例，并链式加载中间件
		s.cron = cron.New(
			cron.WithSeconds(), // 支持 6 位 Cron 表达式（包含秒）
			cron.WithChain(
				cron.Recover(logger),            // 自动捕获 panic，防止单个任务崩掉整个进程
				cron.SkipIfStillRunning(logger), // 如果上一次任务还没执行完，跳过本次触发
			),
			cron.WithLogger(logger),
		)
	})
}

// AddJob 添加带 Cron 表达式的定时任务
func (s *CronService) AddJob(spec string, name string, cmd func()) (cron.EntryID, error) {
	entryID, err := s.cron.AddFunc(spec, func() {
		start := time.Now()
		log.Printf("[Task Start] -> %s", name)
		cmd()
		log.Printf("[Task End]   -> %s (耗时: %v)", name, time.Since(start))
	})

	if err != nil {
		return 0, fmt.Errorf("添加定时任务 [%s] 失败: %w", name, err)
	}
	return entryID, nil
}

// Start 启动定时任务调度器（非阻塞）
func (s *CronService) Start() {
	s.cron.Start()
	log.Println("[CRON] 调度器已启动...")
}

// Stop 优雅关闭调度器，等待正在运行的任务完成
func (s *CronService) Stop() {
	log.Println("[CRON] 正在关闭调度器，等待正在执行的任务结束...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("[CRON] 调度器已安全关闭")
}
