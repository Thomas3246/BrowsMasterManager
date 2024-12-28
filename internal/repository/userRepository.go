package repository

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
)

type UserRepository interface {
	RegisterUser(user *entites.User) error
	ChangeUserName(ctx context.Context, id string, newName string) error
}
