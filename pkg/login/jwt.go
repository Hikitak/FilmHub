package login

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// jwtKey holds the HMAC secret used for signing tokens. It is loaded from the
// JWT_SECRET environment variable. If the variable is not set a hard-coded
// fallback is used, but this should ONLY be relied on during local
// development. In production you **must** provide JWT_SECRET.
var jwtKey = []byte(func() string {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return v
	}
	// TODO: replace with panic or fatal log if running in production without secret
	return "supersecretkey"
}())

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
