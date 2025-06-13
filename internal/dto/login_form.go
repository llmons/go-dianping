package dto

type LoginForm struct {
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
