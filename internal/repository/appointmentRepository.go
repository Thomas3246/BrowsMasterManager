package repository

import "github.com/Thomas3246/BrowsMasterManager/internal/entites"

type AppointmentRepository interface {
	CreateAppointment(appointment *entites.Appointment) error
}
