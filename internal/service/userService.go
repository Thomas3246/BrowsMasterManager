package service

import (
	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepo}
}

func (s *UserService) RegisterUser(id string) error {

	user := entites.User{
		Id:    id,
		Name:  "",
		Phone: "",
	}

	err := s.UserRepository.RegisterUser(&user)
	return err
}
