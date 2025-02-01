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

	query := `INSERT INTO appointments () VALUES ()`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	// for i := 0; i < 10; i++ {
	// 	if ctx.Err() != nil {
	// 		log.Print(ctx.Err())
	// 		return ctx.Err()
	// 	}
	// 	fmt.Println(i)
	// 	time.Sleep(time.Second)
	// }

	return nil
}

func (r *SqliteAppointmentRepository) GetAvailableServices(ctx context.Context) (services []entites.Service, err error) {
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
