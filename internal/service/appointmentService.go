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
)

type AppointmentService struct {
	appointmentRepository repository.AppointmentRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		appointmentRepository: appointmentRepo,
	}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, id int64, appointment *entites.Appointment) error {

	err := s.appointmentRepository.CreateAppointment(ctx, appointment)
	return err
}

func (s *AppointmentService) GetAvailableServices(ctx context.Context) (services []entites.Service, err error) {
	services, err = s.appointmentRepository.GetAvailableServices(ctx)
	if err != nil {
		return nil, err
	}

	return services, err
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
