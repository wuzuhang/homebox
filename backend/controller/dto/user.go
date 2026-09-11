package dto

type RegisterReq struct {
	Id       uint   `json:"id"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`

	Nickname   string  `json:"nickname"`
	Avatar     string  `json:"avatar"`
	Gender     int     `json:"gender"`             // 性别：0:未知, 1:男, 2:女
	Birthday   string  `json:"birthday,omitempty"` // 出生日期（格式：YYYY-MM-DD）
	Weight     float32 `json:"weight"`             // 体重(kg)，辅助儿童用药剂量计算
	DiseaseIDs []uint  `json:"disease_ids"`        // 疾病标签 ID 列表
	Remarks    string  `json:"remarks"`            // 备注
}
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
