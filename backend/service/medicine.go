package service

import (
	"errors"
	"home-box/backend/controller/dto"
	"home-box/backend/models"

	"github.com/jinzhu/copier"
)

type MedicineService struct{}

// 新增药品
func (s *MedicineService) AddMedicine(medicine dto.MedicineReq) error {
	var existingMedicine models.Medicine
	_ = copier.Copy(&existingMedicine, &medicine)
	return models.CreateMedicine(&existingMedicine)
}

// 更新药品
func (s *MedicineService) UpdateMedicine(medicine dto.MedicineReq) error {
	var existingMedicine models.Medicine
	_ = copier.Copy(&existingMedicine, &medicine)
	return models.UpdateMedicine(&existingMedicine)
}

// 删除药品：入参直接接收 uint 类型的 id
func (s *MedicineService) DeleteMedicine(id uint) error {
	if id == 0 {
		return errors.New("无效的药品 ID")
	}
	return models.DeleteMedicine(id)
}

// 查询用户的所有药品
func (s *MedicineService) GetMedicinesByUserID(userID uint) ([]models.Medicine, error) {
	if userID == 0 {
		return nil, errors.New("无效的用户 ID")
	}
	return models.FindMedicineByUserIDs(userID)
}
