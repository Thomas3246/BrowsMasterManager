package sqlite

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type SqliteAppointmentRepository struct {
	db *sql.DB
}

func NewSqliteAppointmentRepository(db *sql.DB) repository.AppointmentRepository {
	return &SqliteAppointmentRepository{db: db}
}

func (r *SqliteAppointmentRepository) CreateAppointment(appointment *entites.Appointment) error {

	// SQL - запрос к БД на INSERT

	return nil
}
