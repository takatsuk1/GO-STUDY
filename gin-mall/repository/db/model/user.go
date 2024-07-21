package model

import (
	"gin-mall/config"
	"gin-mall/pkg/utils/log"
	"github.com/CocaineCong/secret"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Money          string
	Relations      []User `gorm:"many2many:relation"`
}

const (
	PassWordCost        = 12
	Active       string = "active"
)

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}

// 加密金额
func (u *User) EncryptMoney(key string) (money string, err error) {
	aesObj, err := secret.NewAesEncrypt(config.Config.EncryptSecret.MoneySecret, key, "", secret.AesEncrypt128, secret.AesModeTypeCBC)
	if err != nil {
		log.LogrusObj.Error(err)
		return "", err
	}
	money = aesObj.SecretEncrypt(u.Money)
	return
}

// 解密金额
func (u *User) DecryptMoney(key string) (money float64, err error) {
	aesObj, err := secret.NewAesEncrypt(config.Config.EncryptSecret.MoneySecret, key, "", secret.AesEncrypt128, secret.AesModeTypeCBC)
	if err != nil {
		return
	}

	money = cast.ToFloat64(aesObj.SecretDecrypt(u.Money))
	return
}
