package controller

import (
	"home-box/backend/controller/dto"
	"home-box/backend/service"

	"github.com/gin-gonic/gin"
)

var medicineService = &service.MedicineService{}

func AddMedicine(c *gin.Context) {
	var medicine dto.MedicineReq
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数格式错误"})
		return
	}
	if err := medicineService.AddMedicine(medicine); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "添加成功"})
}
func UpdateMedicine(c *gin.Context) {
	var medicine dto.MedicineReq
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数格式错误"})
		return
	}
	if err := medicineService.UpdateMedicine(medicine); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "更新成功"})
}
func DeleteMedicine(c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数格式错误"})
		return
	}
	if err := medicineService.DeleteMedicine(req.ID); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
}
func GetMedicinesByUserID(c *gin.Context) {
	userID := c.GetUint("userID")
	medicines, err := medicineService.GetMedicinesByUserID(userID)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "查询成功", "data": medicines})
}
