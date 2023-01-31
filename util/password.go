package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ハッシュパスワードを生成
// ハッシュパスワードはALG, COST, SALT, HASHで構成される
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash ppassword: %w", err)
	}
	return string(hashPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
