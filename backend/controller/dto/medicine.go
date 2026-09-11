package dto

type MedicineReq struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name" binding:"required"`
	UserID       uint    `json:"user_id"`
	Manufacturer string  `json:"manufacturer"`
	DiseaseIDs   []uint  `json:"disease_ids"`
	Stock        float32 `json:"stock"`
	Unit         string  `json:"unit"`
	MinStockWarn float64 `json:"min_stock_warn"`
	DailyDose    float64 `json:"daily_dose"`
	Usage        string  `json:"usage"`
	Photo        string  `json:"photo"`
	Remark       string  `json:"remark"`
}
