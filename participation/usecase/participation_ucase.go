package usecase

import (
	"context"
	"time"

	"undina/domain"
	"undina/participation"
)

type participationUsecase struct {
	participationRepo domain.ParticipationRepository
	contextTimeout    time.Duration
}

func NewParticipationUsecase(participationRepo domain.ParticipationRepository, timeout time.Duration) domain.ParticipationUsecase {
	return &participationUsecase{
		participationRepo: participationRepo,
		contextTimeout:    timeout,
	}
}

func (p *participationUsecase) Create(c context.Context, participation *domain.Participation) (err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	err = p.participationRepo.Create(ctx, participation)
	if err != nil {
		return err
	}
	return nil
}

func (p *participationUsecase) IsAdmin(c context.Context, roomId int32) (res bool, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := p.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return false, participation.ErrNotHost
	}
	return true, nil
}

func (p *participationUsecase) IsMember(c context.Context, roomId int32) (res bool, err error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isMember, err := p.participationRepo.IsMember(ctx, roomId, user.Id)
	if !isMember || err != nil {
		return false, participation.ErrNotMember
	}
	return true, nil
}
