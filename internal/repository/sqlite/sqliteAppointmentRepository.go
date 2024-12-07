package sqlite

import (
	"context"
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

func (r *SqliteAppointmentRepository) CreateAppointment(ctx context.Context, appointment *entites.Appointment) error {

	// !--------------------!
	//   Make insert query
	// !--------------------!

	query := `INSERT INTO appointments VALUES ()`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
