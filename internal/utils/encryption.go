package utils

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func HashString(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), 14)
	return string(hashed), err
}

func CheckHashAndPassword(hashPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func ReturnSecretFromConfig() string {
	return viper.GetString("secret")
}
