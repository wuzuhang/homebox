package service

import (
	"errors"
	"time"

	"home-box/backend/controller/dto"
	"home-box/backend/models"
	"home-box/utils"

	"github.com/jinzhu/copier"
)

type UserService struct{}

// Register 注册业务逻辑
func (s *UserService) Register(req dto.RegisterReq) error {
	// 1. 检查用户名是否存在
	count := models.CheckUsernameExists(req.Username)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 2. 密码加密（实际项目应使用 bcrypt）
	var user models.User
	var birthdayPtr *time.Time
	if req.Birthday != "" {
		// 手动按照 "2006-01-02" 格式解析
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthdayPtr = &t
		}
	}
	_ = copier.Copy(&user, &req)
	user.Birthday = birthdayPtr

	// 3. 写入数据库
	return models.CreateUser(&user)
}

func (s *UserService) UpdateProfile(userID uint, req dto.RegisterReq) error {
	// 1. 查找用户
	user, err := models.FindUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	// 2. 更新用户信息
	var birthdayPtr *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthdayPtr = &t
		}
	}
	_ = copier.Copy(&user, &req)
	user.Birthday = birthdayPtr

	// 3. 保存到数据库
	return models.UpdateUser(user)
}

// Login 登录业务逻辑，成功则返回 Token
func (s *UserService) Login(username, password string) (string, error) {
	// 1. 查询用户
	user, err := models.FindUserByUsername(username)
	if err != nil {
		return "", errors.New("账号或密码错误")
	}

	// 2. 校验密码
	if user.Password != password {
		return "", errors.New("账号或密码错误")
	}

	// 3. 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", errors.New("生成登录凭证失败")
	}

	return token, nil
}

// GetUserProfile 获取用户信息业务逻辑
func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	return models.FindUserByID(userID)
}
