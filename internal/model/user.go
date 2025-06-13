package model

type User struct {
	Id         uint   `gorm:"primary_key" json:"id"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	NickName   string `json:"nickname"`
	Icon       string `json:"icon"`
	CreateTime string `json:"-"`
	UpdateTime string `json:"-"`
}

func (u *User) TableName() string {
	return "tb_user"
}
