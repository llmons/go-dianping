package model

import "time"

type User struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	NickName   string    `json:"nickname"`
	Icon       string    `json:"icon"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"-"`
	UpdateTime time.Time `gorm:"autoCreateTime" json:"-"`
}

func (u *User) TableName() string {
	return "tb_user"
}
