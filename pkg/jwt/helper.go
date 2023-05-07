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
	Role  string `json:"role"`
}

func GenerateAccessToken(u user.GetUsersDto, role string) (token string, err error) {
	t := jwt.New(jwt.SigningMethodHS256)
	mapClaims := t.Claims.(jwt.MapClaims)
	mapClaims["id"] = u.Id
	mapClaims["email"] = u.Email
	mapClaims["role"] = role
	mapClaims["is_verified"] = u.IsVerified
	mapClaims["exp"] = time.Now().Add(time.Second * 30).Unix()
	token, err = t.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, err
}

func ParseAccessToken(token string) (jwt.MapClaims, error) {
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}

	return claims.Claims.(jwt.MapClaims), err
}

func GenerateRefreshToken(id string) (token string, err error) {
	t := jwt.New(jwt.SigningMethodHS256)
	mapClaims := t.Claims.(jwt.MapClaims)
	mapClaims["id"] = id
	mapClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token, err = t.SignedString([]byte("refresh_token"))
	if err != nil {
		return "", err
	}

	return token, err
}

func ParseRefreshTokenToken(token string) (jwt.MapClaims, error) {
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte("refresh_token"), nil
	})
	if err != nil {
		return nil, err
	}

	return claims.Claims.(jwt.MapClaims), err
}
