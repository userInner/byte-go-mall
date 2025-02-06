package model

type User struct {
	Base
	Email          string `gorm:"unique"`
	PasswordHashed string
}

func (u User) TableName() string {
	return "tb_user"
}
