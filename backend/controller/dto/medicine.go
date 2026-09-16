package dto

type MedicineReq struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name" binding:"required"`
	UserID         uint    `json:"user_id"`
	Manufacturer   string  `json:"manufacturer"`
	DiseaseIDs     []uint  `json:"disease_ids"`
	Specifications uint    `json:"specifications"`
	Price          float32 `json:"price"`
	Stock          float32 `json:"stock"`
	PackageUnit    string  `json:"package_unit"`
	DoseUnit       string  `json:"dose_unit"`
	MinStockWarn   float64 `json:"min_stock_warn"`
	DailyDose      float64 `json:"daily_dose"`
	Usage          string  `json:"usage"`
	Photo          string  `json:"photo"`
	State          int     `json:"state"`
	Remark         string  `json:"remark"`
}
