package auth

import (
    "fmt"
    "time"
    "github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your-secret-key") // In production, use environment variable

type Claims struct {
    UserID string `json:"user_id"`
    jwt.StandardClaims
}

func GenerateToken(userID string) (string, error) {
    claims := Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}