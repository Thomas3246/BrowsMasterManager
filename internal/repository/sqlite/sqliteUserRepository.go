package sqlite

import (
	"context"
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type SqliteUserRepository struct {
	db *sql.DB
}

func NewSqliteUserRepository(db *sql.DB) repository.UserRepository {
	return &SqliteUserRepository{db: db}
}

func (r *SqliteUserRepository) RegisterUser(ctx context.Context, user *entites.User) error {

	query := `INSERT INTO Users (user_id, name, phone_number, role) VALUES (?, ?, ?, (SELECT role_id FROM Roles WHERE role_name = ?))`

	_, err := r.db.ExecContext(ctx, query, user.Id, user.Name, user.Phone, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqliteUserRepository) CheckForUser(ctx context.Context, userId string) (name string, err error) {
	query := `SELECT name FROM Users 
			  WHERE user_id = ?`
	err = r.db.QueryRowContext(ctx, query, userId).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *SqliteUserRepository) ChangeUserName(ctx context.Context, id string, newName string) (err error) {
	query := `UPDATE Users
			  SET name = ?
			  WHERE user_id = ?`

	_, err = r.db.ExecContext(ctx, query, newName, id)
	if err != nil {
		return err
	}

	// for i := 0; i < 12; i++ {
	// 	if ctx.Err() != nil {
	// 		log.Print(ctx.Err())
	// 		return ctx.Err()
	// 	}
	// 	fmt.Println(i)
	// 	time.Sleep(time.Second)
	// }

	return nil
}

func (r *SqliteUserRepository) CheckForAppointments(ctx context.Context, userId int64) (appointments []entites.Appointment, err error) {
	query := `SELECT Appointments.appointment_id, Appointments.date, Appointments.hour, Appointments.minute, Appointments.duration, Appointments.cost
			  FROM Appointments 
			  WHERE Appointments.user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		appointment := entites.Appointment{}
		err = rows.Scan(&appointment.ID, &appointment.DateStr, &appointment.Hour, &appointment.Minute, &appointment.TotalDuration, &appointment.TotalCost)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *SqliteUserRepository) CheckForUserByPhone(ctx context.Context, phone string) (id int, err error) {
	query := "SELECT user_id FROM Users WHERE phone_number = ?"
	row := r.db.QueryRowContext(ctx, query, phone)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
