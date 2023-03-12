package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func GenerateAccessToken(u user.GetUsersDto) (token string, err error) {
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)),
		},
		Email: u.Email,
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	token, err = claims.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return token, err
}

func ParseAccessToken(token string) (jwt.Claims, error) {
	claims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}

	return claims.Claims, err
}
