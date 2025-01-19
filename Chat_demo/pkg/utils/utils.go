package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Format_Uid(uid1, uid2 string) string {
	return fmt.Sprintf("p%sSEP%sp", uid1, uid2)
}

func HashCompare(hash_str, src_str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash_str), []byte(src_str))
	return err == nil
}

func HashStr(str string) (string, error) {
	PassWordCost := 12                                                   // 加密难度
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), PassWordCost) // 调用加密算法
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 自定义错误类型
type MyError struct {
	Message string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("MyError: %s", e.Message)
}
