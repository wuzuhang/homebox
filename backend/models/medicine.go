package models

import (
	"home-box/config"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Medicine 药品实体表
type Medicine struct {
	gorm.Model

	UserID uint `gorm:"not null;index" json:"user_id"` // 所属用户 ID

	// 1. 基础属性
	Name         string `gorm:"type:varchar(100);not null;index" json:"name"` // 药品商品名/通用名（如：“布洛芬缓释胶囊”）
	Manufacturer string `gorm:"type:varchar(100)" json:"manufacturer"`        // 生产企业/药厂

	// 2. 存放与分类
	// 👈 直接以 JSON 形式存储疾病标签 ID 列表，例如：[1, 3, 5]
	DiseaseIDs datatypes.JSONSlice[uint] `gorm:"type:json" json:"disease_ids,omitempty"`

	// 3. 库存与规格
	Stock        float32 `gorm:"type:decimal(8,2);default:0" json:"stock"`          // 当前剩余库存数量
	Unit         string  `gorm:"type:varchar(20);default:'粒'" json:"unit"`          // 包装/规格单位（粒 / 片 / 盒 / 支 / ml）
	MinStockWarn float64 `gorm:"type:decimal(8,2);default:2" json:"min_stock_warn"` // 低库存预警阈值
	DailyDose    float64 `gorm:"type:decimal(8,2);default:0" json:"daily_dose"`     // 每天用药剂量

	// 5. 用法用量与禁忌
	Usage string `gorm:"type:varchar(255)" json:"usage"` // 用法用量说明（如：“口服，一次1粒，一日2次，饭后服用”）

	// 7. 附加信息
	Photo  string `gorm:"type:varchar(255)" json:"photo"`  // 药品外观/包装照片 URL
	Remark string `gorm:"type:varchar(255)" json:"remark"` // 备注
}

func CreateMedicine(medicine *Medicine) error {
	return config.DB.Create(medicine).Error
}
func DeleteMedicine(id uint) error {
	return config.DB.Delete(id).Error
}
func UpdateMedicine(medicine *Medicine) error {
	return config.DB.Model(&Medicine{}).
		Where("id = ?", medicine.ID).
		Updates(medicine).Error
}
func FindMedicineByUserIDs(id uint) ([]Medicine, error) {
	var medicines []Medicine
	err := config.DB.Where("user_id = ?", id).Find(&medicines).Error
	if err != nil {
		return nil, err
	}
	return medicines, nil
}
