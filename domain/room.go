package domain

import (
	"context"
)

type RoomStatus string

const (
	CREATED RoomStatus = "created"
	START   RoomStatus = "start"
	END     RoomStatus = "end"
)

type PaymentStatus string

const (
	UNPAID    PaymentStatus = "unpaid"
	PENDING   PaymentStatus = "pending"
	CONFIRMED PaymentStatus = "confirmed"
)

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

type RoomRepository interface {
	GetJoinedRooms(ctx context.Context, id int32) ([]RoomItem, error)
}

type RoomUsecase interface {
	GetJoinedRooms(ctx context.Context) ([]RoomItem, error)
}
