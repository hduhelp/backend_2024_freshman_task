package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	CustomClaims map[string]interface{}
}

func CreateJWT(secretKey string, ttl time.Duration, claims map[string]interface{}) (string, error) {
	expirationTime := time.Now().Add(ttl)
	jwtClaims := jwt.MapClaims{
		"exp": expirationTime.Unix(),
	}
	for key, value := range claims {
		jwtClaims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(secretKey))
}

func ParseJWT(secretKey, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
