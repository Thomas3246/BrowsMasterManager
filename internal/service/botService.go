package service

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/repository/sqlite"
)

type BotService struct {
	UserService        *UserService
	AppointmentService *AppointmentService
}

func NewBotService(db *sql.DB) *BotService {

	appointmentRepo := sqlite.NewSqliteAppointmentRepository(db)
	appointmentService := NewAppointmentService(appointmentRepo)

	userRepo := sqlite.NewSqliteUserRepository(db)
	userService := NewUserService(userRepo)

	return &BotService{
		UserService:        userService,
		AppointmentService: appointmentService,
	}
}
