package usecase

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"time"

	"undina/domain"
	// adapterqueue "undina/internal/adapter/queue"
	// helper "undina/internal/common"
	// "undina/internal/infrastructure/queue"
	"undina/room"

	"github.com/sirupsen/logrus"
)

type roomUsecase struct {
	roomRepo          domain.RoomRepository
	participationRepo domain.ParticipationRepository
	serviceRepo       domain.ServiceRepository
	invitationRepo    domain.InvitationRepository
	roundRepo         domain.RoundRepository
	contextTimeout    time.Duration
	// producer          adapterqueue.Producer
}

// func NewRoomUsecase(roomRepo domain.RoomRepository, participationRepo domain.ParticipationRepository, serviceRepo domain.ServiceRepository, invitationRepo domain.InvitationRepository, roundRepo domain.RoundRepository, timeout time.Duration, queue *queue.RabbitMQHandler) domain.RoomUsecase {
func NewRoomUsecase(roomRepo domain.RoomRepository, participationRepo domain.ParticipationRepository, serviceRepo domain.ServiceRepository, invitationRepo domain.InvitationRepository, roundRepo domain.RoundRepository, timeout time.Duration) domain.RoomUsecase {
	return &roomUsecase{
		roomRepo:          roomRepo,
		participationRepo: participationRepo,
		serviceRepo:       serviceRepo,
		invitationRepo:    invitationRepo,
		roundRepo:         roundRepo,
		contextTimeout:    timeout,
		// 	producer:          adapterqueue.NewProducer(queue.Channel(), "paymentCheck"),
	}
}

func (r *roomUsecase) Create(c context.Context, roomRequest *domain.RoomRequest) (res int32, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	roomRequest.AdminId = user.Id

	plan, err := r.serviceRepo.GetPlanByKey(ctx, roomRequest.PlanName, fmt.Sprintf("%d", roomRequest.ServiceId))
	if err != nil {
		return 0, err
	}

	if plan.MaxCount < roomRequest.MaxCount {
		return 0, room.ErrMaxCountExceed
	}

	room := &domain.Room{
		ServiceId: roomRequest.ServiceId,
		PlanName:  roomRequest.PlanName,
		MaxCount:  roomRequest.MaxCount,
		AdminId:   roomRequest.AdminId,
		IsPublic:  roomRequest.IsPublic,
	}

	if *room.IsPublic {
		room.PublicMessage = roomRequest.PublicMessage
		room.MatchingDeadline = sql.NullString{String: roomRequest.MatchingDeadline, Valid: true}
		room.RoomStatus = domain.CREATED
	} else {
		room.RoomStatus = domain.START
	}

	roomId, err := r.roomRepo.Create(ctx, room)
	if err != nil {
		return 0, err
	}

	err = r.participationRepo.Create(ctx, &domain.Participation{
		UserId:        user.Id,
		RoomId:        roomId,
		PaymentStatus: domain.CONFIRMED,
		IsHost:        true,
	})
	if err != nil {
		return 0, err
	}

	return roomId, nil
}

func (r *roomUsecase) Start(c context.Context, roomId int32) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return room.ErrNotHost
	}

	roomInfo, err := r.GetRoomInfo(ctx, roomId)
	if err != nil {
		return err
	}
	if *roomInfo.RoomStatus != domain.CREATED {
		return room.ErrAlreadyStarted
	}

	err = r.roomRepo.Start(ctx, roomId)
	if err != nil {
		return err
	}

	return nil
}

func (r *roomUsecase) GetPublicRooms(c context.Context) (res []domain.RoomPublic, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	res, err = r.roomRepo.GetPublicRooms(ctx)
	if err != nil {
		return nil, err
	}

	for i, room := range res {
		if room.MatchingDeadline != "" {
			matchingDeadline, _ := time.Parse(time.RFC3339, room.MatchingDeadline)
			res[i].MatchingDeadline = fmt.Sprintf("%d/%02d/%02d", matchingDeadline.Year(), matchingDeadline.Month(), matchingDeadline.Day())
		}
	}

	return res, nil
}

func (r *roomUsecase) GetJoinedRooms(c context.Context) (res []domain.RoomItem, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	res, err = r.participationRepo.GetJoinedRooms(ctx, user.Id)
	for i, roomItem := range res {
		itemSplitCost, err := r.GetRoomSplitFee(ctx, roomItem.RoomId)
		if err != nil {
			continue
		}
		res[i].Cost = itemSplitCost
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roomUsecase) GenerateInvitationCode(c context.Context, roomId int32) (res string, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return "", room.ErrNotHost
	}

	code := sha1.New()
	code.Write([]byte(time.Now().String()))
	code.Write([]byte(fmt.Sprint(roomId)))
	invitationCode := fmt.Sprintf("%x", code.Sum(nil))[0:7]

	err = r.invitationRepo.GenerateInvitationCode(ctx, roomId, invitationCode)
	if err != nil {
		return "", err
	}

	return invitationCode, nil
}

func (r *roomUsecase) GetInvitationCodes(c context.Context, roomId int32) (res []domain.InvitationCode, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return nil, room.ErrNotHost
	}

	res, err = r.invitationRepo.GetInvitationCodes(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roomUsecase) JoinRoom(c context.Context, code string) (res int32, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	roomId, err := r.invitationRepo.ConsumeInvitationCode(ctx, code)
	if err != nil {
		return 0, room.ErrInvalidInvitationCode
	}

	roomInfo, err := r.participationRepo.GetRoomInfo(ctx, roomId)
	members, err := r.participationRepo.GetRoomMembers(ctx, roomId)

	if len(members) >= int(roomInfo.MaxCount) {
		r.invitationRepo.ResumeInvitationCode(ctx, code)
		return 0, room.ErrRoomFull
	}

	user := c.Value(domain.CtxUserKey).(*domain.User)

	err = r.participationRepo.Create(ctx, &domain.Participation{
		UserId:        user.Id,
		RoomId:        roomId,
		PaymentStatus: domain.UNPAID,
		IsHost:        false,
	})
	if err != nil {
		r.invitationRepo.ResumeInvitationCode(ctx, code)
		return 0, room.ErrAlreadyJoined
	}

	return roomId, nil
}

func (r *roomUsecase) LeaveRoom(c context.Context, roomId int32, userId int32) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return room.ErrNotHost
	}

	err = r.participationRepo.LeaveRoom(ctx, roomId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *roomUsecase) GetRoomInfo(c context.Context, roomId int32) (res *domain.RoomInfoResponse, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	_, err = r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if err != nil {
		return nil, room.ErrNotMember
	}

	res, err = r.participationRepo.GetRoomInfo(ctx, roomId)
	if err != nil {
		return nil, err
	}

	if res.MatchingDeadline != "" {
		matchingDeadline, _ := time.Parse(time.RFC3339, res.MatchingDeadline)
		res.MatchingDeadline = fmt.Sprintf("%d/%02d/%02d", matchingDeadline.Year(), matchingDeadline.Month(), matchingDeadline.Day())
	}

	return res, nil
}

func (r *roomUsecase) GetRoomAdmin(c context.Context, roomId int32) (res *domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	_, err = r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if err != nil {
		return nil, room.ErrNotMember
	}

	res, err = r.participationRepo.GetRoomAdmin(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roomUsecase) GetRoomMembers(c context.Context, roomId int32) (res []domain.Participation, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	_, err = r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if err != nil {
		return nil, room.ErrNotMember
	}

	res, err = r.participationRepo.GetRoomMembers(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *roomUsecase) GetRoomSplitFee(c context.Context, roomId int32) (res int32, err error) {
	members, err := r.participationRepo.GetRoomMembers(c, roomId)
	if err != nil {
		return 0, err
	}
	roomFeeInfo, err := r.participationRepo.GetRoomFeeInfo(c, roomId)
	if err != nil {
		return 0, err
	}
	var avgFee = (roomFeeInfo.Cost * roomFeeInfo.RoundInterval) / int32(len(members))
	return avgFee, nil
}

func (r *roomUsecase) UpdateRoom(c context.Context, roomId int32, roomRequest *domain.RoomRequest) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	_, err = r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if err != nil {
		return room.ErrNotHost
	}

	plan, err := r.serviceRepo.GetPlanByKey(ctx, roomRequest.PlanName, fmt.Sprintf("%d", roomRequest.ServiceId))
	if err != nil {
		return err
	}

	if plan.MaxCount < roomRequest.MaxCount {
		return room.ErrMaxCountExceed
	}

	err = r.roomRepo.Update(ctx, roomId, &domain.Room{
		ServiceId:    roomRequest.ServiceId,
		PlanName:     roomRequest.PlanName,
		MaxCount:     roomRequest.MaxCount,
		IsPublic:     roomRequest.IsPublic,
		Announcement: roomRequest.Announcement,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *roomUsecase) Delete(c context.Context, roomId int32) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return room.ErrNotHost
	}

	err = r.roomRepo.Delete(ctx, roomId)
	if err != nil {
		return err
	}

	return nil
}

func (r *roomUsecase) GetTodayStartingMember(c context.Context) (res []domain.ParticipationInfo, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	// truncate timestamp to date only
	now := time.Now().Truncate(24 * time.Hour)
	res, err = r.participationRepo.GetRoomMemberByStartingTime(ctx, now)
	if err != nil {
		return nil, err
	}
	for i, participantInfo := range res {
		// get average fee of the room
		owed_fee, err := r.GetRoomSplitFee(c, participantInfo.RoomId)
		if err != nil {
			return nil, err
		}
		res[i].OwedFee = owed_fee
		// get admin Info
		adminRes, err := r.participationRepo.GetRoomAdmin(c, participantInfo.RoomId)
		if err != nil {
			return nil, err
		}
		res[i].AdminName = adminRes.Name
		res[i].AdminEmail = adminRes.Email
	}

	return res, nil
}

func (r *roomUsecase) GetTodayPaymentDueMember(c context.Context) (res []domain.ParticipationInfo, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	// truncate timestamp to date only
	now := time.Now().Truncate(24 * time.Hour)
	res, err = r.participationRepo.GetRoomMemberByDueTime(ctx, now)
	if err != nil {
		return nil, err
	}
	for i, participantInfo := range res {
		// get average fee of the room
		owed_fee, err := r.GetRoomSplitFee(c, participantInfo.RoomId)
		if err != nil {
			return nil, err
		}
		res[i].OwedFee = owed_fee
		// get admin Info
		adminRes, err := r.participationRepo.GetRoomAdmin(c, participantInfo.RoomId)
		if err != nil {
			return nil, err
		}
		res[i].AdminName = adminRes.Name
		res[i].AdminEmail = adminRes.Email
	}

	return res, nil
}

func (r *roomUsecase) GetParticipationInfoByRoomId(c context.Context, roomId int32) (res []domain.ParticipationInfo, err error) {
	res, err = r.participationRepo.GetRoomMemberById(c, roomId)
	if err != nil {
		return nil, err
	}
	logrus.Info(res)

	for i, participantInfo := range res {
		// get average fee of the room
		owed_fee, err := r.GetRoomSplitFee(c, participantInfo.RoomId)
		if err != nil {
			return nil, err
		}
		res[i].OwedFee = owed_fee
		// get admin Info
		adminRes, err := r.participationRepo.GetRoomAdmin(c, participantInfo.RoomId)
		if err != nil {
			return nil, err
		}
		res[i].AdminName = adminRes.Name
		res[i].AdminEmail = adminRes.Email
	}
	return res, nil
}

func (r *roomUsecase) AddRound(c context.Context, roomId int32, roundRequest *domain.RoundRequest) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return room.ErrNotHost
	}

	roomInfo, err := r.GetRoomInfo(ctx, roomId)
	if *roomInfo.RoomStatus == domain.CREATED {
		return room.ErrNotStarted
	}

	roundInfo, err := r.GetRound(ctx, roomId)
	if roundInfo.StartingTime != "" {
		return room.ErrRoundAlreadyCreated
	}

	start, err := time.Parse("2006-01-02", roundRequest.StartingTime)
	if err != nil {
		logrus.Info("parse time err: ", err)
		return err
	}
	end := start.AddDate(0, int(roundRequest.RoundInterval), 0)
	deadline := start.AddDate(0, 0, -(roundRequest.PaymentDeadline * 7))

	roundId, err := r.roundRepo.AddRound(ctx, &domain.Round{
		StartingTime:    start,
		EndingTime:      end,
		RoundInterval:   roundRequest.RoundInterval,
		PaymentDeadline: deadline,
		IsAddCalendar:   roundRequest.IsAddCalendar,
	})
	if err != nil {
		return err
	}

	err = r.roomRepo.UpdateRoundId(ctx, roomId, roundId)
	if err != nil {
		return err
	}

	// send email to user in rooms to inform start new round
	// participationInfos, err := r.GetParticipationInfoByRoomId(c, roomId)
	// for _, info := range participationInfos {
	// 	message := helper.EncodeToBytes(&info)
	// 	message = helper.Compress(message)
	// 	r.producer.Publish(message)
	// }
	return nil
}

func (r *roomUsecase) GetRound(c context.Context, roomId int32) (res *domain.RoundInfo, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	res, err = r.roundRepo.GetRound(ctx, roomId)
	if res.StartingTime != "" {
		start, _ := time.Parse(time.RFC3339, res.StartingTime)
		end, _ := time.Parse(time.RFC3339, res.EndingTime)
		deadline, _ := time.Parse(time.RFC3339, res.PaymentDeadline)

		res.StartingTime = fmt.Sprintf("%d/%02d/%02d", start.Year(), start.Month(), start.Day())
		res.EndingTime = fmt.Sprintf("%d/%02d/%02d", end.Year(), end.Month(), end.Day())
		res.PaymentDeadline = fmt.Sprintf("%d/%02d/%02d", deadline.Year(), deadline.Month(), deadline.Day())
	}

	return res, nil
}

func (r *roomUsecase) DeleteRound(c context.Context, roomId int32) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	if !isAdmin || err != nil {
		return room.ErrNotHost
	}

	roomInfo, err := r.GetRoomInfo(ctx, roomId)
	if err != nil {
		return err
	}

	if roomInfo.RoundId != 0 {
		err := r.roundRepo.DeleteRound(ctx, roomInfo.RoundId)
		if err != nil {
			return err
		}
	} else {
		return room.ErrNoRound
	}
	return nil
}

func (r *roomUsecase) UpdatePaymentStatus(c context.Context, roomId int32, userId int32, paymentStatus domain.PaymentStatus) (err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	user := c.Value(domain.CtxUserKey).(*domain.User)

	isAdmin, err := r.participationRepo.IsAdmin(ctx, roomId, user.Id)
	// only user himself or admin can change payment status
	if (!isAdmin && user.Id != userId) || err != nil {
		return room.ErrNotAuthorized
	}

	err = r.participationRepo.UpdatePaymentStatus(ctx, roomId, userId, paymentStatus)
	if err != nil {
		return err
	}

	return nil
}

func (r *roomUsecase) IsPublic(c context.Context, roomId int32) (res bool, err error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	res, err = r.roomRepo.IsPublic(ctx, roomId)
	if err != nil {
		return res, err
	}
	return res, nil
}
