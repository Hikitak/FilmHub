package login

import (
	"errors"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	secretOnce sync.Once
	jwtKey     []byte
)

// Init sets the HMAC secret used for signing JWTs. It must be called once at
// application startup. Re-initialization is ignored.
func Init(secret string) {
	secretOnce.Do(func() {
		jwtKey = []byte(secret)
	})
}

func prepareKey() error {
	if len(jwtKey) == 0 {
		return errors.New("jwt secret not initialized: call login.Init first")
	}
	return nil
}

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, role string) (string, error) {
	if err := prepareKey(); err != nil {
		return "", err
	}
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
	if err := prepareKey(); err != nil {
		return nil, err
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
