package repository

import "github.com/Thomas3246/BrowsMasterManager/internal/entites"

type UserRepository interface {
	RegisterUser(user *entites.User) error
}
