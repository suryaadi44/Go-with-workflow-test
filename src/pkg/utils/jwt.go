package utils

import (
	"rewrite/pkg/config"
	"rewrite/pkg/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    user.ID,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT_SECRET))
}
