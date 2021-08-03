package usecase

import (
	"context"
	"time"

	"undina/domain"
)

type serviceUsecase struct {
	serviceRepo    domain.ServiceRepository
	contextTimeout time.Duration
}

func NewServiceUsecase(serviceRepo domain.ServiceRepository, timeout time.Duration) domain.ServiceUsecase {
	return &serviceUsecase{
		serviceRepo:    serviceRepo,
		contextTimeout: timeout,
	}
}

func (s *serviceUsecase) FetchAll(c context.Context) (res []domain.Service, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	res, err = s.serviceRepo.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUsecase) GetDetailByID(c context.Context, id string) (res []domain.ServiceDetail, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	res, err = s.serviceRepo.GetDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUsecase) GetPlanByKey(c context.Context, planName string, serviceId string) (res *domain.Plan, err error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	res, err = s.serviceRepo.GetPlanByKey(ctx, planName, serviceId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
