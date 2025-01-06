package service

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, id string, phone string) error {

	user := entites.User{
		Id:    id,
		Name:  "",
		Phone: "",
	}

	err := s.UserRepository.RegisterUser(&user)
	return err
}

func (s *UserService) CheckForUser(ctx context.Context, phone string) (name string) {
	// вызов репо, возврат имени пользователя в случае, если зареган.  "" - Если нет
	return "name"
}

func (s *UserService) ChangeUserName(ctx context.Context, id string, newName string) (err error) {

	err = s.UserRepository.ChangeUserName(ctx, id, newName)
	return err
}
