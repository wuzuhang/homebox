package controller

import (
	"home-box/backend/controller/dto"
	"home-box/config"

	"github.com/gin-gonic/gin"
)

func GetDiseases(c *gin.Context) {
	// 这里可以调用 TagService 来获取标签列表
	tags := []dto.Diseases{}
	err := config.DB.Find(&tags).Error

	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "获取标签失败"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "获取标签成功", "data": tags})
}
