package controllers

import (
	"backend3/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User) (string, error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})
	token, err := generateToken.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return token, err
	}
	return token, nil
}
