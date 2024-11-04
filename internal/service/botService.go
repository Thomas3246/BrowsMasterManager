package service

import (
	"log"

	"github.com/Thomas3246/BrowsMasterManager/internal/repository/sqlite"
)

type BotService struct {
	RegisterUserService *RegisterUserService
	BookingService      *BookingService
}

func NewBotService() *BotService {

	//To-DO мб перенести отсюда InitDB

	// Сервис не должен знать про БД

	db, err := sqlite.InitDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// TO-DO сделать по-другому объявление сервисов и их добавление

	bookingRepo := sqlite.NewSqliteBookingRepository(db)

	bookingService := NewBookingService(bookingRepo)

	return &BotService{
		RegisterUserService: nil,
		BookingService:      bookingService,
	}
}
