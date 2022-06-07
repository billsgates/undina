package domain

import (
	"context"
)

type Application struct {
	UserId             int32   `json:"user_id,omitempty"`
	UserName           string  `json:"user_name,omitempty"`
	UserRating         float32 `json:"user_rating"`
	ApplicationMessage string  `json:"application_message"`
	ApplicationDate    string  `json:"application_date,omitempty"`
	RoomId             int32   `json:"room_id,omitempty"`
	IsAccepted         bool    `json:"is_accepted"`
}

type ApplicationRequest struct {
	UserId             int32  `json:"user_id,omitempty"`
	RoomId             int32  `json:"room_id,omitempty"`
	ApplicationMessage string `json:"application_message,omitempty"`
}

type ApplicationUsecase interface {
	Create(ctx context.Context, roomId int32, message string) error
	FetchAll(ctx context.Context, roomId int32) (res []Application, err error)
	IsApplied(ctx context.Context, roomId int32) (res bool, err error)
	AcceptApplication(ctx context.Context, roomId int32, userId int32) (err error)
	DeleteApplication(ctx context.Context, roomId int32, userId int32) (err error)
}

type ApplicationRepository interface {
	Create(ctx context.Context, applicationRequest *ApplicationRequest) error
	FetchAll(ctx context.Context, roomId int32) (res []Application, err error)
	IsApplied(ctx context.Context, roomId int32, userId int32) (res bool, err error)
	AcceptApplication(ctx context.Context, roomId int32, userId int32) (err error)
	DeleteApplication(ctx context.Context, roomId int32, userId int32) (err error)
}
