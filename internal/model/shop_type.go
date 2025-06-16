package model

import "time"

type ShopType struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Sort       uint      `json:"sort"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

func (m *ShopType) TableName() string {
	return "tb_shop_type"
}
