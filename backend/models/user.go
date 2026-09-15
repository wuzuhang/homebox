package models

import (
	"home-box/config"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User 系统登录账号表
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;type:varchar(50);not null" json:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"-"`       // 密码不上报 JSON
	Phone    string `gorm:"uniqueIndex;type:varchar(20)" json:"phone"` // 手机号（可用于验证码登录/通知）
	Email    string `gorm:"type:varchar(100)" json:"email"`

	Nickname string     `gorm:"type:varchar(50)" json:"nickname"`
	Avatar   string     `gorm:"type:varchar(255)" json:"avatar"`      // 头像
	Gender   int        `gorm:"type:tinyint;default:0" json:"gender"` // 性别：0:未知, 1:男, 2:女
	Birthday *time.Time `json:"birthday,omitempty"`                   // 出生日期（格式：YYYY-MM-DD）
	Weight   float32    `gorm:"type:decimal(5,2)" json:"weight"`      // 体重(kg)，辅助儿童用药剂量计算
	Remarks  string     `gorm:"type:varchar(255)" json:"remarks"`     // 备注

	// 👈 直接以 JSON 形式存储疾病标签 ID 列表，例如：[1, 3, 5]
	DiseaseIDs datatypes.JSONSlice[uint] `gorm:"type:json" json:"disease_ids,omitempty"`
}

// CheckUsernameExists 校验用户名是否存在
func CheckUsernameExists(username string) int64 {
	var count int64
	config.DB.Model(&User{}).Where("username = ?", username).Count(&count)
	return count
}

// FindUserByUsername 根据用户名查找用户
func FindUserByUsername(username string) (*User, error) {
	var user User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// CreateUser 创建新用户
func CreateUser(user *User) error {
	return config.DB.Create(user).Error
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	// 指定 Model(user) 后，GORM 会自动识别 user.ID 作为 WHERE 条件
	return config.DB.Model(user).Updates(user).Error
}

// FindUserByID 根据 ID 查找用户
func FindUserByID(id uint) (*User, error) {
	var user User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
