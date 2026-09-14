package cron

import (
	"log"

	"home-box/backend/models"
	"home-box/config"

	"gorm.io/gorm"
)

func UpdateMedicineStock() {
	// 使用表达式进行原子扣减，并将小于 0 的结果限制为 0
	// GREATEST(0, stock - daily_dose) 可以防止库存变成负数
	err := config.DB.Model(&models.Medicine{}).
		Where("state = ? AND daily_dose > 0 AND stock > 0", 1).
		Update("stock", gorm.Expr("GREATEST(0, stock - daily_dose)")).Error

	if err != nil {
		log.Printf("[Cron Error] 扣减药品每日用量失败: %v\n", err)
		return
	}

	log.Println("[Cron Info] 药品每日用量扣减完成")
}
