package usecase

import (
	"context"
	"fmt"
	"time"

	"undina/domain"
)

type applicationUsecase struct {
	applicationRepo domain.ApplicationRepository
	contextTimeout  time.Duration
}

func NewApplicationUsecase(applicationRepo domain.ApplicationRepository, timeout time.Duration) domain.ApplicationUsecase {
	return &applicationUsecase{
		applicationRepo: applicationRepo,
		contextTimeout:  timeout,
	}
}

func (a *applicationUsecase) Create(c context.Context, roomId int32, message string) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	// err = a.applicationRepo.Create(ctx, roomId, user.Id)
	err = a.applicationRepo.Create(ctx, &domain.ApplicationRequest{
		RoomId:             roomId,
		UserId:             user.Id,
		ApplicationMessage: message,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *applicationUsecase) FetchAll(c context.Context, roomId int32) (res []domain.Application, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.applicationRepo.FetchAll(ctx, roomId)
	if err != nil {
		return nil, err
	}

	for i, application := range res {
		applicationTime, _ := time.Parse(time.RFC3339, application.ApplicationDate)
		res[i].ApplicationDate = fmt.Sprintf("%d/%02d/%02d", applicationTime.Year(), applicationTime.Month(), applicationTime.Day())
	}

	return res, nil
}

func (a *applicationUsecase) IsApplied(c context.Context, roomId int32) (res bool, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	res, err = a.applicationRepo.IsApplied(ctx, roomId, user.Id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (a *applicationUsecase) AcceptApplication(c context.Context, roomId int32, userId int32) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err = a.applicationRepo.AcceptApplication(ctx, roomId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *applicationUsecase) DeleteApplication(c context.Context, roomId int32, userId int32) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err = a.applicationRepo.DeleteApplication(ctx, roomId, userId)
	if err != nil {
		return err
	}
	return nil
}
