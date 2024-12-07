package repository

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
)

type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, appointment *entites.Appointment) error
}
