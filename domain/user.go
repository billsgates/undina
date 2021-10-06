package domain

import (
	"context"
)

type User struct {
	Id             int32   `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Email          string  `json:"email,omitempty"`
	PasswordDigest string  `json:"password_digest,omitempty"`
	Phone          string  `json:"phone"`
	Rating         float32 `json:"rating"`
	RatingCount    int32   `json:"rating_count"`
}

type UserRequest struct {
	Id       int32  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmailPassword(ctx context.Context, email string, password string) (*User, error)
}

type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmailPassword(ctx context.Context, email string, password string) (*User, error)
}
