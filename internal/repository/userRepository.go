package repository

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user *entites.User) error
	CheckForUser(ctx context.Context, userId string) (string, error)
	CheckForAppointments(ctx context.Context, userId int64) (appointments []entites.Appointment, err error)
	ChangeUserName(ctx context.Context, id string, newName string) error
	CheckForUserByPhone(ctx context.Context, phone string) (id int, err error)
}
