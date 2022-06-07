package domain

import (
	"context"
	"database/sql"
	"time"
)

type RoomStatus string

const (
	CREATED RoomStatus = "created"
	START   RoomStatus = "start"
	END     RoomStatus = "end"
)

type Room struct {
	Id               int32          `json:"id,omitempty"`
	Announcement     string         `json:"announcement,omitempty"`
	IsPublic         *bool          `json:"is_public,omitempty"`
	RoomStatus       RoomStatus     `json:"room_status,omitempty"`
	CreatedAt        time.Time      `json:"created_at,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty"`
	MaxCount         int32          `json:"max_count,omitempty"`
	AdminId          int32          `json:"admin_id,omitempty"`
	ServiceId        int32          `json:"service_id,omitempty"`
	PlanName         string         `json:"plan_name,omitempty"`
	PublicMessage    string         `json:"public_message,omitempty"`
	MatchingDeadline sql.NullString `json:"matching_deadline,omitempty"`
}

type RoomRequest struct {
	MaxCount         int32  `json:"max_count" binding:"required"`
	AdminId          int32  `json:"admin_id,omitempty"`
	ServiceId        int32  `json:"service_id" binding:"required"`
	PlanName         string `json:"plan_name" binding:"required"`
	PaymentPeriod    int32  `json:"payment_period"`
	IsPublic         *bool  `json:"is_public" binding:"required"`
	Announcement     string `json:"announcement,omitempty"`
	PublicMessage    string `json:"public_message,omitempty"`
	MatchingDeadline string `json:"matching_deadline,omitempty"`
}

type RoomJoinRequest struct {
	InvitationCode string `json:"invitation_code" binding:"required"`
}

type RoomItem struct {
	RoomId        int32          `json:"room_id"`
	ServiceName   string         `json:"service_name"`
	PlanName      string         `json:"plan_name"`
	IsHost        bool           `json:"is_host"`
	PaymentStatus *PaymentStatus `json:"payment_status"`
	RoomStatus    *RoomStatus    `json:"room_status"`
	Cost          int32          `json:"cost,omitempty"`
	IsPublic      bool           `json:"is_public"`
}

type RoomInfoResponse struct {
	RoomId           int32           `json:"room_id,omitempty"`
	IsPublic         *bool           `json:"is_public"`
	Announcement     string          `json:"announcement"`
	MatchingDeadline string          `json:"matching_deadline"`
	MaxCount         int32           `json:"max_count,omitempty"`
	RoomStatus       *RoomStatus     `json:"room_status,omitempty"`
	RoundId          int32           `json:"round_id"`
	ServiceId        int32           `json:"service_id,omitempty"`
	ServiceName      string          `json:"service_name,omitempty"`
	PlanName         string          `json:"plan_name,omitempty"`
	Role             string          `json:"role,omitempty"`
	PaymentFee       int32           `json:"payment_fee,omitempty"`
	Round            *RoundInfo      `json:"round" gorm:"-"`
	Admin            *User           `json:"admin,omitempty" gorm:"-"`
	Members          []Participation `json:"members,omitempty" gorm:"-"`
}

type RoomPublic struct {
	RoomId           int32   `json:"room_id,omitempty"`
	AdminName        string  `json:"admin_name,omitempty"`
	AdminRating      float32 `json:"admin_rating"`
	ServiceName      string  `json:"service_name,omitempty"`
	PlanName         string  `json:"plan_name,omitempty"`
	MaxCount         int32   `json:"max_count,omitempty"`
	MemberCount      int32   `json:"member_count,omitempty"`
	Cost             int32   `json:"cost,omitempty"`
	MatchingDeadline string  `json:"matching_deadline,omitempty"`
	PublicMessage    string  `json:"public_message,omitempty"`
	IsApplied        bool    `json:"is_applied"`
}

type RoomFeeInfo struct {
	RoomId        int32  `json:"room_id"`
	ServiceName   string `json:"service_name"`
	PlanName      string `json:"plan_name"`
	Cost          int32  `json:"cost"`
	RoundInterval int32  `json:"round_interval,omitempty"`
}

type RoomRepository interface {
	Create(ctx context.Context, room *Room) (roomId int32, err error)
	GetPublicRooms(ctx context.Context) (res []RoomPublic, err error)
	Update(ctx context.Context, roomId int32, room *Room) error
	UpdateRoundId(ctx context.Context, roomId int32, roundId int32) error
	Delete(ctx context.Context, roomId int32) (err error)
	Start(ctx context.Context, roomId int32) (err error)
	IsPublic(ctx context.Context, roomId int32) (res bool, err error)
}

type RoomUsecase interface {
	Create(ctx context.Context, room *RoomRequest) (roomId int32, err error)
	Delete(ctx context.Context, roomId int32) error
	Start(ctx context.Context, roomId int32) error
	GetRoomInfo(ctx context.Context, roomId int32) (res *RoomInfoResponse, err error)
	GetRoomAdmin(ctx context.Context, roomId int32) (res *User, err error)
	GetRoomMembers(ctx context.Context, roomId int32) (res []Participation, err error)
	GetRoomSplitFee(ctx context.Context, roomId int32) (res int32, err error)
	GetPublicRooms(ctx context.Context) (res []RoomPublic, err error)
	GetTodayStartingMember(c context.Context) (res []ParticipationInfo, err error)
	GetTodayPaymentDueMember(c context.Context) (res []ParticipationInfo, err error)
	GetParticipationInfoByRoomId(c context.Context, roomId int32) (res []ParticipationInfo, err error)
	GetJoinedRooms(ctx context.Context) ([]RoomItem, error)
	GenerateInvitationCode(ctx context.Context, roomId int32) (string, error)
	GetInvitationCodes(ctx context.Context, roomId int32) (res []InvitationCode, err error)
	JoinRoom(ctx context.Context, code string) (roomId int32, err error)
	LeaveRoom(ctx context.Context, roomId int32, userId int32) error
	UpdateRoom(ctx context.Context, roomId int32, room *RoomRequest) error
	UpdatePaymentStatus(ctx context.Context, userId int32, roomId int32, status PaymentStatus) error
	AddRound(ctx context.Context, roomId int32, round *RoundRequest) error
	GetRound(ctx context.Context, roomId int32) (res *RoundInfo, err error)
	DeleteRound(ctx context.Context, roomId int32) error
	IsPublic(ctx context.Context, roomId int32) (res bool, err error)
}
