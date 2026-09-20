package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// ServerConfig 服务相关配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库相关配置
type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type COSConfig struct {
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
}

// Config 总配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"` // 👈 增加此行
	COS      COSConfig      `mapstructure:"cos"` // 👈 增加此行
}

// 全局配置变量
var AppConfig *Config

// InitConfig 初始化并加载配置文件
func InitConfig() {
	v := viper.New()
	v.SetConfigFile("config.yaml") // 配置文件文件名
	v.SetConfigType("yaml")        // 配置文件类型

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 将配置内容反序列化到 AppConfig 结构体中
	AppConfig = &Config{}
	if err := v.Unmarshal(AppConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	fmt.Println("✅ 配置文件加载成功！")
}
