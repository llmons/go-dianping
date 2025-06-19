package model

type User struct {
	Model
	Phone    string `gorm:"unique"`
	Password string
	NickName string
	Icon     string
}

func (u *User) TableName() string {
	return "tb_user"
}
