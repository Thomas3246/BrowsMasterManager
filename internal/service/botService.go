package service

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/repository/sqlite"
)

type BotService struct {
	UserService        *UserService
	AppointmentService *AppointmentService
	ServiceService     *ServiceService
}

func NewBotService(db *sql.DB) *BotService {

	appointmentRepo := sqlite.NewSqliteAppointmentRepository(db)
	appointmentService := NewAppointmentService(appointmentRepo)

	userRepo := sqlite.NewSqliteUserRepository(db)
	userService := NewUserService(userRepo)

	serviceRepo := sqlite.NewSqliteServiceRepository(db)
	serviceService := NewServiceService(serviceRepo)

	return &BotService{
		UserService:        userService,
		AppointmentService: appointmentService,
		ServiceService:     serviceService,
	}
}
