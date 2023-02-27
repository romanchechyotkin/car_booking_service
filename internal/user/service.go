package user

import (
	"context"
	user2 "github.com/romanchechyotkin/car_booking-service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking-service/internal/user/storage"
)

type Service struct {
	repository *user.Repository
}

func NewService(repository *user.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll() ([]user2.GetUsersDto, error) {
	users, err := s.repository.GetAllUsers(context.Background())
	return users, err
}

func (s *Service) CreateUser(ctx context.Context, user *user2.CreateUserDto) error {
	err := s.repository.CreateUser(ctx, user)
	return err
}
