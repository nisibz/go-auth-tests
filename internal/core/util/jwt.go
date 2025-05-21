package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nisibz/go-auth-tests/internal/adapter/config"
)

var jwtSecretKey []byte

const tokenExpirationDuration = 24 * time.Hour // Token valid for 24 hours

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func InitJWTSecretKey(cfg *config.Container) error {
	if cfg == nil || cfg.JwtSecretKey == nil || cfg.JwtSecretKey.JWT_SECRET_KEY == "" {
		return fmt.Errorf("JWT secret key not found in config")
	}
	jwtSecretKey = []byte(cfg.JwtSecretKey.JWT_SECRET_KEY)
	return nil
}

func GenerateToken(userID string) (string, error) {
	if len(jwtSecretKey) == 0 {
		return "", fmt.Errorf("JWT_SECRET_KEY environment variable not set or empty")
	}

	expirationTime := time.Now().Add(tokenExpirationDuration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	if len(jwtSecretKey) == 0 {
		return "", fmt.Errorf("JWT_SECRET_KEY environment variable not set or empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
