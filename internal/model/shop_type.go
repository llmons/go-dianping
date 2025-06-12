package model

import "time"

type ShopType struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Sort       uint      `json:"sort"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

func (m *ShopType) TableName() string {
	return "tb_shop_type"
}
