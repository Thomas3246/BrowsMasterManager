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

func (r *SqliteUserRepository) RegisterUser(user *entites.User) error {

	// add ctx

	// user CHECK-OUT

	// INSERT UserData query

	return nil
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
