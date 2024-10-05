package Token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

var jwtKey = []byte("<random>")

// GenerateToken 新建JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 令牌过期时间
		Id:        strconv.Itoa(int(userID)),             //转为字符串并设置为Id
		Subject:   username,                              //设置用户名
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //使用jwt.NewWithClaims方法，结合签名方法jwt.SigningMethodHS256和声明信息
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil //调用token.SignedString(jwtKey)方法，利用之前定义的密钥对令牌进行签名，并返回签名后的字符串
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	//接收待解析的JWT令牌字符串
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证令牌的签名方法是否为预期的HS256
		if alg := token.Method.Alg(); alg != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", alg)
		}
		//调用jwt.Parse方法，尝试解析传入的令牌字符串
		//第二个参数是一个回调函数，用于验证令牌的签名。在此函数中，我们返回之前定义的jwtKey作为验证密钥
		return jwtKey, nil //jwt.Parse方法返回一个*jwt.Token类型的值和一个error。如果解析成功且签名有效，error将为nil
	})
	if err != nil {
		// 如果解析或签名验证失败，返回错误
		return nil, err
	}

	// 检查令牌是否有效
	if !token.Valid {
		return nil, errors.New("invalid token: token is not valid")
	}

	// 使用类型断言检查token.Claims是否为*jwt.StandardClaims类型
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		// 如果类型断言失败，返回错误
		return nil, fmt.Errorf("claims verification failed")
	}

	return claims, nil
}
