package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secret-key")

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return signedToken, nil

}

func main() {
	tokenStr, err := createToken("john_doe")
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	fmt.Println("JWT Token:", tokenStr)
	token, _ := ValidateJWT(tokenStr)
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println("Username:", claims["username"])
		fmt.Println("Token expires at:", time.Unix(int64(claims["exp"].(float64)), 0))
	} else {
		fmt.Println("Invalid token")
	}

}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
