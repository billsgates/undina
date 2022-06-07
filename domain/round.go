package domain

import (
	"context"
	"time"
)

type RoundRequest struct {
	RoomId          int32  `json:"room_id,omitempty"`
	StartingTime    string `json:"starting_time,omitempty" binding:"required"`
	RoundInterval   int32  `json:"round_interval,omitempty" binding:"required"`
	PaymentDeadline int    `json:"payment_deadline,omitempty" binding:"required"`
	IsAddCalendar   *bool  `json:"is_add_calendar,omitempty" binding:"required"`
}

type Round struct {
	RoundId         int32     `json:"round_id,omitempty" gorm:"primary_key"`
	StartingTime    time.Time `json:"starting_time,omitempty"`
	EndingTime      time.Time `json:"ending_time,omitempty"`
	RoundInterval   int32     `json:"round_interval,omitempty"`
	PaymentDeadline time.Time `json:"payment_deadline,omitempty"`
	IsAddCalendar   *bool     `json:"is_add_calendar,omitempty"`
}

type RoundInfo struct {
	StartingTime    string `json:"starting_time"`
	EndingTime      string `json:"ending_time"`
	PaymentDeadline string `json:"payment_deadline"`
	RoundInterval   int32  `json:"round_interval"`
}

type RoundRepository interface {
	AddRound(ctx context.Context, round *Round) (roundId int32, err error)
	GetRound(ctx context.Context, roomId int32) (res *RoundInfo, err error)
	DeleteRound(ctx context.Context, roundId int32) error
}
