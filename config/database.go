package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// 使用 AppConfig 中读取的驱动和 DSN
	switch AppConfig.Database.Driver {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(AppConfig.Database.Dsn), &gorm.Config{})
	default:
		log.Fatalf("不支持的数据库驱动: %s", AppConfig.Database.Driver)
	}

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	fmt.Println("✅ 数据库初始化成功！")
}
