package model

import "time"

type ShopType struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Sort       uint      `json:"sort"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (m *ShopType) TableName() string {
	return "tb_shop_type"
}
