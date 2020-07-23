package password

import (
	"github.com/kataras/golog"
	"golang.org/x/crypto/bcrypt"
)

func PasswordToHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	encodePW := ""
	if err != nil {
		golog.Default.Debug(err)
		return encodePW, err
	} else {
		// 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
		encodePW = string(hash)
	}
	return encodePW, err
}

func PasswordCheck(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
