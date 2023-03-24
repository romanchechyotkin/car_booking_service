package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"

	auth "github.com/romanchechyotkin/car_booking_service/internal/auth/model"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"

	"fmt"
	"log"
)

type service struct {
	repository *user.Repository
}

func NewService(rep *user.Repository) *service {
	return &service{
		repository: rep,
	}
}

func (s *service) Registration(ctx *gin.Context, dto auth.RegistrationDto) error {
	_, userErr := s.repository.GetOneUserByEmail(ctx, dto.Email)
	if userErr == nil {
		return fmt.Errorf("email is used")
	}

	password := dto.Password
	hashedPassword, _ := hashPassword(password)

	var cu = user2.CreateUserDto{
		Email:           dto.Email,
		Password:        hashedPassword,
		FullName:        dto.FullName,
		TelephoneNumber: dto.TelephoneNumber,
	}

	createErr := s.repository.CreateUser(ctx, &cu)
	if createErr != nil {
		return fmt.Errorf("telephone number is used")
	}

	log.Printf("user %s registrated", cu.Email)
	return nil
}

func (s *service) Login(ctx *gin.Context, dto auth.LoginDto) (u user2.GetUsersDto, token string, err error) {
	u, userErr := s.repository.GetOneUserByEmail(ctx, dto.Email)
	if userErr != nil {
		return u, "", fmt.Errorf("user not found")
	}

	hashedPassword := checkPasswordHash(dto.Password, u.Password)
	if !hashedPassword {
		return u, "", fmt.Errorf("wrong password")
	}

	token, err = jwt.GenerateAccessToken(u)
	if err != nil {
		return u, "", err
	}

	log.Printf("user %s logined", u.Email)
	return u, token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}