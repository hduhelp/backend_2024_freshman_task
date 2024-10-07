package auth

import (
	"QA/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("tobeempireofhduhelp")

// 加密密码
func GenerateSalt() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic("Error generating salt")
	}
	return hex.EncodeToString(bytes)
}

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func HashPassword(password, salt string) string {
	saltedPassword := password + salt
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))
	return hex.EncodeToString(hash.Sum(nil))
}

// jwt
func GenerateJWT(u models.User) (string, error) {
	claims := MyClaims{
		Username: u.Name,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,         // Token生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // Token过期时间
			Issuer:    u.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 使用HS256算法
	tokenString, err := token.SignedString(jwtKey)             // 使用密钥进行签名
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func AuthenticateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Println(claims.Username)
		return token, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
