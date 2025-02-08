package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
	rusdate "github.com/Thomas3246/BrowsMasterManager/pkg/rusDate"
)

type SqliteAppointmentRepository struct {
	db *sql.DB
}

func NewSqliteAppointmentRepository(db *sql.DB) repository.AppointmentRepository {
	return &SqliteAppointmentRepository{db: db}
}

func (r *SqliteAppointmentRepository) CreateAppointment(ctx context.Context, appointment *entites.Appointment) error {

	query := `INSERT INTO Appointments (user_id, date, hour, minute, duration) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, appointment.UserId, rusdate.FormatDayMonth(appointment.Date), appointment.Hour, appointment.Minute, appointment.TotalDuration)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	getIDQuery := "SELECT appointment_id FROM Appointments WHERE user_id = ? AND date = ? AND hour = ? AND minute = ?"
	appId := r.db.QueryRowContext(ctx, getIDQuery, appointment.UserId, rusdate.FormatDayMonth(appointment.Date), appointment.Hour, appointment.Minute)
	err = appId.Scan(&appointment.ID)
	if err != nil {
		log.Println("Ошибка определения ID записи: ", err)
		return err
	}

	appointmentServiceQuery := fmt.Sprintf("INSERT INTO appointments_services (appointment_id, service_id) VALUES (%d, ?)", appointment.ID)
	for i := range appointment.Services {
		if appointment.Services[i].Added {
			_, err = r.db.ExecContext(ctx, appointmentServiceQuery, appointment.Services[i].Id)
			if err != nil {
				log.Println("Ошибка добавления записи в таблицу appointments_services: ", err)
				return err
			}
		}
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

func (r *SqliteAppointmentRepository) CheckAppointmentsAtDate(ctx context.Context, date string) (appointments []entites.Appointment, err error) {
	query := "SELECT hour, minute, duration FROM Appointments WHERE date = ?"

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		appointment := entites.Appointment{}
		err = rows.Scan(&appointment.Hour, &appointment.Minute, &appointment.TotalDuration)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}
