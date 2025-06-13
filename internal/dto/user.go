package dto

type User struct {
	Id       uint   `json:"id"`
	NickName string `json:"nickname"`
	Icon     string `json:"icon"`
}
