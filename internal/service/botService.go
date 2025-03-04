package service

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/repository/sqlite"
	"github.com/redis/go-redis/v9"
)

type BotService struct {
	UserService        *UserService
	AppointmentService *AppointmentService
	ServiceService     *ServiceService
}

func NewBotService(db *sql.DB, redis *redis.Client) *BotService {

	appointmentRepo := sqlite.NewSqliteAppointmentRepository(db)
	appointmentService := NewAppointmentService(appointmentRepo, redis)

	userRepo := sqlite.NewSqliteUserRepository(db)
	userService := NewUserService(userRepo, redis)

	serviceRepo := sqlite.NewSqliteServiceRepository(db)
	serviceService := NewServiceService(serviceRepo, redis)

	return &BotService{
		UserService:        userService,
		AppointmentService: appointmentService,
		ServiceService:     serviceService,
	}
}
