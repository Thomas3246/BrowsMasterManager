package sqlite

import (
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

	// user CHECK-OUT

	// INSERT UserData query

	return nil
}
