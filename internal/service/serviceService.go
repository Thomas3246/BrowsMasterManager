package service

import (
	"context"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository"
)

type ServiceService struct {
	serviceRepository repository.ServiceRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository) *ServiceService {
	return &ServiceService{
		serviceRepository: serviceRepo,
	}
}

func (s *ServiceService) GetAvailableServices(ctx context.Context) (services []entites.Service, err error) {
	services, err = s.serviceRepository.GetAvailableServices(ctx)
	if err != nil {
		return nil, err
	}

	return services, err
}

func (s *ServiceService) GetServicesInAppointment(ctx context.Context, id int) (services []entites.Service, err error) {
	services, err = s.serviceRepository.GetServicesInAppointment(ctx, id)
	if err != nil {
		return nil, err
	}

	return services, nil
}
