package auth

import (
	"QA/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

func GenerateSalt() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic("Error generating salt")
	}
	return hex.EncodeToString(bytes)
}

func HashPassword(password, salt string) string {
	saltedPassword := password + salt
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateJWT(u models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{
		"id":   u.StudentID,
		"name": u.Name,
	}
	token.Claims = claims
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthenticateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return token, fmt.Errorf("invalid token")
	}
	return token, nil
}
