package sqlite

import (
	"context"
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type SqliteServiceRepository struct {
	db *sql.DB
}

func NewSqliteServiceRepository(db *sql.DB) repository.ServiceRepository {
	return &SqliteServiceRepository{db: db}
}

func (r *SqliteServiceRepository) GetAvailableServices(ctx context.Context) (services []entites.Service, err error) {
	query := `SELECT * FROM Services`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		service := entites.Service{}
		err = rows.Scan(&service.Id, &service.Name, &service.Descr, &service.Cost, &service.Duration)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, err
}

func (r *SqliteServiceRepository) GetServicesInAppointment(ctx context.Context, appointmentId int) (services []entites.Service, err error) {
	query := `SELECT Appointments_services.service_id, Services.name
			  FROM Services INNER JOIN Appointments_services ON appointments_services.service_id = Services.service_id
			  WHERE appointments_services.appointment_id = ?`

	rows, err := r.db.QueryContext(ctx, query, appointmentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		service := entites.Service{}
		err = rows.Scan(&service.Id, &service.Name)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}
