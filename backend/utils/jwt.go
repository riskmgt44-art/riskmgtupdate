// utils/jwt.go
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"riskmgt/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserID primitive.ObjectID `json:"userID"`
	Name   string             `json:"name"`
	Role   string             `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID primitive.ObjectID, role, name string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour * 7)

	claims := &Claims{
		UserID: userID,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}