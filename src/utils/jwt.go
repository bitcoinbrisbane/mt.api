package utils

import (
	"os"
	"time"

	"example.com/mt/src/models"
	"github.com/golang-jwt/jwt"
)

// GenerateToken generates a new JWT token for a given username
func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &models.JWTClaims{
		Username: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	var jwtSecret = os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

// ValidateToken validates the JWT token and returns the user claims
func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}
	var jwtSecret = os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err //jwt.NewValidationError("Invalid token")
	}

	return claims, nil
}