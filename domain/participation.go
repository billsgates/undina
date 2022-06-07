package domain

import (
	"context"
	"time"
)

type PaymentStatus string

const (
	UNPAID    PaymentStatus = "unpaid"
	PENDING   PaymentStatus = "pending"
	CONFIRMED PaymentStatus = "confirmed"
)

type Participation struct {
	UserId        int32         `json:"user_id,omitempty"`
	UserName      string        `json:"user_name,omitempty"`
	RoomId        int32         `json:"room_id,omitempty"`
	PaymentStatus PaymentStatus `json:"payment_status,omitempty"`
	IsHost        bool          `json:"is_host,omitempty"`
}

type ParticipationInfo struct {
	UserId          int32  `json:"user_id,omitempty"`
	UserName        string `json:"user_name,omitempty"`
	UserEmail       string `json:"user_email,omitempty"`
	ServiceProvider string `json:"service_provider,omitempty"`
	PlanName        string `json:"plan_name,omitempty"`
	OwedFee         int32  `json:"owed_fee,omitempty"`
	RoomId          int32  `json:"room_id,omitempty"`
	AdminId         int32  `json:"admin_id,omitempty"`
	AdminName       string `json:"admin_name,omitempty"`
	AdminEmail      string `json:"admin_email,omitempty"`
}

type ParticipationRequest struct {
	UserId int32 `json:"user_id,omitempty" binding:"required"`
	RoomId int32 `json:"room_id,omitempty" binding:"required"`
}

type ParticipationStatusRequest struct {
	UserId        int32         `json:"user_id,omitempty" binding:"required"`
	RoomId        int32         `json:"room_id,omitempty" binding:"required"`
	PaymentStatus PaymentStatus `json:"payment_status,omitempty" binding:"required"`
}

type ParticipationUsecase interface {
	Create(ctx context.Context, participation *Participation) error
	IsMember(ctx context.Context, roomId int32) (bool, error)
	IsAdmin(ctx context.Context, roomId int32) (bool, error)
}

type ParticipationRepository interface {
	Create(ctx context.Context, participation *Participation) error
	GetRoomInfo(ctx context.Context, roomId int32) (res *RoomInfoResponse, err error)
	GetRoomAdmin(ctx context.Context, roomId int32) (res *User, err error)
	GetRoomMembers(ctx context.Context, roomId int32) (res []Participation, err error)
	GetJoinedRooms(ctx context.Context, userId int32) ([]RoomItem, error)
	GetRoomFeeInfo(ctx context.Context, roomId int32) (res *RoomFeeInfo, err error)
	GetRoomMemberByStartingTime(ctx context.Context, starting_time time.Time) (res []ParticipationInfo, err error)
	GetRoomMemberByDueTime(ctx context.Context, due_time time.Time) (res []ParticipationInfo, err error)
	GetRoomMemberById(ctx context.Context, roomId int32) (res []ParticipationInfo, err error)
	UpdatePaymentStatus(ctx context.Context, userId int32, roomId int32, status PaymentStatus) error
	IsAdmin(ctx context.Context, roomId int32, userId int32) (bool, error)
	IsMember(ctx context.Context, roomId int32, userId int32) (bool, error)
	LeaveRoom(ctx context.Context, roomId int32, userId int32) error
}
