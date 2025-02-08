package repository

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
)

type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, appointment *entites.Appointment) error
	GetAvailableServices(ctx context.Context) ([]entites.Service, error)
	CheckAppointmentsAtDate(ctx context.Context, date string) (appointments []entites.Appointment, err error)
}
