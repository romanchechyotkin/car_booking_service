package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	user "github.com/romanchechyotkin/car_booking_service/internal/user/model"

	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}

func GenerateAccessToken(u user.GetUsersDto, role string) (token string, err error) {
	uc := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			Subject:   role,
		},
		Email: u.Email,
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	token, err = claims.SignedString([]byte("secret"))
	if err != nil {
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

func GenerateRefreshToken(id string) (token string, err error) {
	uc := jwt.RegisteredClaims{
		Issuer:    id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	token, err = claims.SignedString([]byte("refresh_token"))
	if err != nil {
		return "", err
	}

	return token, err
}

func ParseRefreshTokenToken(token string) (jwt.Claims, error) {
	claims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("refresh_token"), nil
	})
	if err != nil {
		return nil, err
	}

	return claims.Claims, err
}
