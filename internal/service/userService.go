package service

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
	rusdate "github.com/Thomas3246/BrowsMasterManager/pkg/rusDate"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	UserRepository repository.UserRepository
	RedisClient    *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client) *UserService {
	return &UserService{UserRepository: userRepo, RedisClient: redisClient}
}

func (s *UserService) RegisterUser(ctx context.Context, id string, phoneNumber string) error {

	userRole := "client"

	masterPhone, err := s.GetMasterPhone()
	if err != nil {
		return err
	}

	if phoneNumber == masterPhone {
		userRole = "master"
	}

	user := entites.User{
		Id:    id,
		Name:  "",
		Phone: phoneNumber,
		Role:  userRole,
	}

	err = s.UserRepository.RegisterUser(ctx, &user)
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

func (s *UserService) CheckForAppointments(ctx context.Context, userId int64) (validAppointments []entites.Appointment, err error) {
	allAppointments, err := s.UserRepository.CheckForAppointments(ctx, userId)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	for _, appointment := range allAppointments {

		appointmentDate, err := rusdate.FormatBack(appointment.DateStr)
		if err != nil {
			log.Printf("Ошибка обратного форматирования даты: %v", err)
			return nil, err
		}

		hour, _ := strconv.Atoi(appointment.Minute)
		minute, _ := strconv.Atoi(appointment.Minute)
		appointmentTime := appointmentDate.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)

		if currentTime.Before(appointmentTime) {
			validAppointments = append(validAppointments, appointment)
		}
	}

	return validAppointments, nil
}

func (s *UserService) SetMasterPhone(phone string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.RedisClient.Set(ctx, "masterPhone", phone, 0).Err()
	if err != nil {
		log.Printf("Произошла ошибка добавления ключа masterPhone: %s", err)
		return err
	}
	return nil
}

func (s *UserService) GetMasterPhone() (phone string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	phone, err = s.RedisClient.Get(ctx, "masterPhone").Result()
	if err != nil {
		log.Printf("Произошла ошибка получения ключа masterPhone: %s", err)
		return "", err
	}

	return phone, nil
}
