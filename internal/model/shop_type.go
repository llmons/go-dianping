package model

type ShopType struct {
	Model
	Name string
	Icon string
	Sort uint
}

func (m *ShopType) TableName() string {
	return "tb_shop_type"
}
