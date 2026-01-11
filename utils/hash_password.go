package utils

import (
	"golang.org/x/crypto/bcrypt"
)

//pembuatan hash password ketika user di buat amaka akan membuat hash password secara otomatis agar menjaga keamanan akun

func HashPassword(password string) string {
	//pembuatan password hash 14 karakter
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword)
}

func CompareHashAndPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}
