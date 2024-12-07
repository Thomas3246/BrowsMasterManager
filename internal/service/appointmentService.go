package service

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type AppointmentService struct {
	appointmentRepository repository.AppointmentRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		appointmentRepository: appointmentRepo,
	}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, id int64) error {

	appointment := entites.Appointment{
		ID: id,
	}

	err := s.appointmentRepository.CreateAppointment(ctx, &appointment)
	return err
}
