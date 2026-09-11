package dto

type Diseases struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}
