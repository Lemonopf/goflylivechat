package tools

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

func MakeToken(obj map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(obj))
	tokenString, err := token.SignedString([]byte(GetServerConf().JwtSecret))
	return tokenString, err
}
func ParseToken(tokenStr string) map[string]interface{} {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(GetServerConf().JwtSecret), nil
	})
	if err != nil {
		return nil
	}
	finToken, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	return finToken
}
