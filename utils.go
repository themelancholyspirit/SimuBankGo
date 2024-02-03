package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func isValidEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func IsValidRegisterRequest(request *RegisterAccountRequest) bool {

	arr := []string{
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
	}

	for _, el := range arr {
		if el == "" {
			return false
		}
	}

	return isValidEmail(request.Email)
}

func IsValidLoginRequest(request *LoginAccountRequest) bool {
	return isValidEmail(request.Email) && request.Password != ""
}

func CreateToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Invalid token")
	}
}
