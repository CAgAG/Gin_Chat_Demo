package model

import (
	"Chat_demo/pkg/utils"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	hash_password, err := utils.HashStr(password)
	if err != nil {
		return err
	}
	user.Password = hash_password
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return utils.HashCompare(user.Password, password)
}
