package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
	rusdate "github.com/Thomas3246/BrowsMasterManager/pkg/rusDate"
	"github.com/redis/go-redis/v9"
)

type AppointmentService struct {
	appointmentRepository repository.AppointmentRepository
	redisClient           *redis.Client
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, redisClient *redis.Client) *AppointmentService {
	return &AppointmentService{
		appointmentRepository: appointmentRepo,
		redisClient:           redisClient,
	}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, id int64, appointment *entites.Appointment) error {

	err := s.appointmentRepository.CreateAppointment(ctx, appointment)
	return err
}

func (s *AppointmentService) CheckAppointmentsAtDate(ctx context.Context, appointment *entites.Appointment) (appointmentsAtDate []entites.Appointment, err error) {
	appointmentsAtDate, err = s.appointmentRepository.CheckAppointmentsAtDate(ctx, rusdate.FormatDayMonth(appointment.Date))
	if err != nil && err != sql.ErrNoRows {
		log.Println("Ошибка проверки записей на дату: ", err)
		return nil, err
	}

	for i := range appointmentsAtDate {
		hour, _ := strconv.Atoi(appointmentsAtDate[i].Hour)
		minute, _ := strconv.Atoi(appointmentsAtDate[i].Minute)
		appointmentsAtDate[i].Date = appointment.Date.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
	}

	return appointmentsAtDate, nil
}

func (s *AppointmentService) CheckIsBusy(appointmentsAtDate []entites.Appointment, startTime time.Time, totalDuration int) bool {
	endTime := startTime.Add(time.Duration(totalDuration) * time.Minute)

	for _, appointment := range appointmentsAtDate {
		startBusy := appointment.Date
		endBusy := startBusy.Add(time.Duration(appointment.TotalDuration) * time.Minute)

		equals := startTime.Equal(startBusy) || startTime.Equal(endBusy) || endTime.Equal(startBusy) || endTime.Equal(endBusy)
		if equals || (startTime.Before(startBusy) && endTime.After(startBusy)) || (startTime.Before(endBusy) && endTime.After(endBusy)) || (startTime.After(startBusy) && endTime.Before(endBusy)) || (startTime.Before(startBusy) && endTime.After(endBusy)) {
			return true
		}
	}
	return false
}

func (s *AppointmentService) SetAppointmentsInCash(ctx context.Context, id int64, appointments []entites.Appointment) error {
	userId := strconv.Itoa(int(id))

	data, err := json.Marshal(appointments)
	if err != nil {
		log.Printf("Произошла ошибка сериализации: %v", err)
		return err
	}

	err = s.redisClient.Set(ctx, userId, data, 0).Err()
	if err != nil {
		log.Printf("Ошибка записи в redis: %v", err)
		return err
	}
	return nil
}

func (s *AppointmentService) GetAppointmentsFromCash(ctx context.Context, userId int) (appointments []entites.Appointment, err error) {
	id := strconv.Itoa(userId)
	data, err := s.redisClient.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		log.Printf("Ошибка получения записей из redis: %v", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &appointments)
	if err != nil {
		log.Printf("Ошибка десериализации записей из JSON: %v", err)
		return nil, err
	}

	return appointments, nil
}

func (s *AppointmentService) CancelAppointment(appointmentId string, userId int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := strconv.Atoi(appointmentId)
	if err != nil {
		log.Printf("Произошла ошибка преобразования id записи в int: %v", err)
		return err
	}

	err = s.appointmentRepository.CancelAppointment(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Printf("Произошла ошибка при отмене записи: %v", err)
		return err
	}

	userAppointments, err := s.GetAppointmentsFromCash(ctx, int(userId))
	if err != nil {
		return err
	}
	var newUserAppointments []entites.Appointment
	for _, appointment := range userAppointments {
		if appointment.ID != id {
			newUserAppointments = append(newUserAppointments, appointment)
		}
	}

	err = s.SetAppointmentsInCash(ctx, userId, newUserAppointments)
	if err != nil {
		return err
	}

	return nil
}
