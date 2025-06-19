package model

import "time"

type Model struct {
	Id         uint      `gorm:"primary_key"`
	CreateTime time.Time `gorm:"autoCreateTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime"`
}
