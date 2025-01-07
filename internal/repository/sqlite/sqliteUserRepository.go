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

	query := `INSERT INTO Users (user_id, name, phone_number) VALUES (?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, user.Id, user.Name, user.Phone)
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
			  WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query, newName, id)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
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
