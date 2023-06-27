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
	//placer     *emailproducer.EmailPlacer
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
		City:            dto.City,
	}

	createErr := s.repository.CreateUser(ctx, &cu)
	if createErr != nil {
		log.Printf("Error: %v", createErr)
		return fmt.Errorf(createErr.Error())
	}

	//err := s.placer.SendEmail(cu.Email, "registration")
	//if err != nil {
	//	log.Printf("error due kafka %v\n", err)
	//}

	log.Printf("user %s registrated", cu.Email)
	return nil
}

func (s *service) Login(ctx *gin.Context, dto auth.LoginDto) (u user2.GetUsersDto, role, token string, err error) {
	u, userErr := s.repository.GetOneUserByEmail(ctx, dto.Email)
	if userErr != nil {
		return u, "", "", fmt.Errorf("user not found")
	}

	hashedPassword := checkPasswordHash(dto.Password, u.Password)
	if !hashedPassword {
		return u, "", "", fmt.Errorf("wrong password")
	}

	role, err = s.repository.GetRole(ctx, u.Id)
	if err != nil {
		return u, "", "", err
	}
	token, err = jwt.GenerateAccessToken(u, role)
	if err != nil {
		return u, "", "", err
	}

	log.Printf("user %s logined", u.Email)
	return u, role, token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
