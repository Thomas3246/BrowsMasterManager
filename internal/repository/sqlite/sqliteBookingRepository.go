package sqlite

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type SqliteBookingRepository struct {
	db *sql.DB
}

func NewSqliteBookingRepository(db *sql.DB) repository.BookingRepository {
	return &SqliteBookingRepository{db: db}
}

func (d *SqliteBookingRepository) CreateAppointment(appointment *entites.Appointment) error {

	// SQL - запрос к БД на INSERT

	return nil
}
