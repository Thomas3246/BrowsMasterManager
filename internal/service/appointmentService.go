package service

import (
	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type BookingService struct {
	bookingRepository repository.BookingRepository
}

func NewBookingService(bookingRepo repository.BookingRepository) *BookingService {
	return &BookingService{
		bookingRepository: bookingRepo,
	}
}

func (s *BookingService) CreateAppointment(id int64) error {

	appointment := entites.Appointment{
		ID: id,
	}

	err := s.bookingRepository.CreateAppointment(&appointment)
	return err
}
