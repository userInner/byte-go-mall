package model

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	Username string `gorm:"username"`
	Email    string `gorm:"unique"`
	Password string
}

func (u *User) TableName() string {
	return "tb_user"
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
