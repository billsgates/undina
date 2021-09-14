package domain

import (
	"context"

	"github.com/dgrijalva/jwt-go/v4"
)

const CtxUserKey = "user"

type AuthClaims struct {
	jwt.StandardClaims
	User *User `json:"user"`
}

type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Id    int32  `json:"id"`
	Token string `json:"token"`
}

type AuthUsecase interface {
	SignUp(ctx context.Context, name string, email string, password string) (res *LoginResponse, err error)
	SignIn(ctx context.Context, email string, password string) (res *LoginResponse, err error)
	ParseToken(ctx context.Context, accessToken string) (*User, error)
}
