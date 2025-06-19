package model

type Shop struct {
	Model
	Name      string `gorm:"not null"`
	TypeId    int    `gorm:"unique"`
	Images    string `gorm:"not null"`
	Area      string
	Address   string  `gorm:"not null"`
	X         float64 `gorm:"not null"`
	Y         float64 `gorm:"not null"`
	AvgPrice  int
	Sold      int    `gorm:"not null"`
	Comments  string `gorm:"not null"`
	Score     int    `gorm:"not null"`
	OpenHours string
}

func (m *Shop) TableName() string {
	return "tb_shop"
}
