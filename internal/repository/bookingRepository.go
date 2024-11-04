package repository

import "github.com/Thomas3246/BrowsMasterManager/internal/entites"

type BookingRepository interface {
	CreateAppointment(appointment *entites.Appointment) error
}
