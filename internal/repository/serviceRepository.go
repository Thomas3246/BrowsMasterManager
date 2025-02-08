package repository

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
)

type ServiceRepository interface {
	GetAvailableServices(ctx context.Context) ([]entites.Service, error)
	GetServicesInAppointment(ctx context.Context, appointmentId int) (services []entites.Service, err error)
}
