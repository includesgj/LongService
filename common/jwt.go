package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"time"
)

func CreateJWT(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["id"] = id

	// 签名
	tokenString, err := token.SignedString([]byte("智者不入爱河"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return []byte("智者不入爱河"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("用户:", claims["id"])
		stuId, _ := claims["id"].(string)
		return stuId, nil
	} else {
		return "", fmt.Errorf("无效的 token")
	}
}
