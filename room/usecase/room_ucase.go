package usecase

import (
	"context"
	"time"

	"undina/domain"
)

type roomUsecase struct {
	roomRepo       domain.RoomRepository
	contextTimeout time.Duration
}

func NewRoomUsecase(roomRepo domain.RoomRepository, timeout time.Duration) domain.RoomUsecase {
	return &roomUsecase{
		roomRepo:       roomRepo,
		contextTimeout: timeout,
	}
}

func (r *roomUsecase) GetJoinedRooms(c context.Context) (res []domain.RoomItem, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	res, err = r.roomRepo.GetJoinedRooms(ctx, user.Id)
	// for i, roomItem := range res {
	// 	itemSplitCost, err := r.GetRoomSplitFee(ctx, roomItem.RoomId)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	res[i].Cost = itemSplitCost
	// }
	if err != nil {
		return nil, err
	}

	return res, nil
}
