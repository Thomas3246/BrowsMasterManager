package service

import (
	"context"
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, id string, phoneNumber string) error {

	user := entites.User{
		Id:    id,
		Name:  "",
		Phone: phoneNumber,
	}

	err := s.UserRepository.RegisterUser(ctx, &user)
	return err
}

func (s *UserService) CheckForUser(ctx context.Context, userId string) (name string, isRegistred bool, err error) {

	name, err = s.UserRepository.CheckForUser(ctx, userId)
	if err != nil {
		if err != sql.ErrNoRows {
			// Если ошибка не в том, что пользователь не найден, то возвращается ошибка
			return "", isRegistred, err
		}
		// Если ошибка в том, что пользователь не найден
		return "", isRegistred, nil
	}

	isRegistred = true
	return name, isRegistred, nil
}

func (s *UserService) ChangeUserName(ctx context.Context, id string, newName string) (err error) {

	err = s.UserRepository.ChangeUserName(ctx, id, newName)
	return err
}

func (s *UserService) CheckForAppointments(ctx context.Context, userId int64) (appointments []entites.Appointment, err error) {
	appointments, err = s.UserRepository.CheckForAppointments(ctx, userId)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}
