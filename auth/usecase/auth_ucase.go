package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"undina/auth"
	"undina/domain"

	"github.com/dgrijalva/jwt-go/v4"
)

type authUsecase struct {
	userUsecase    domain.UserUsecase
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUsecase(
	userUsecase domain.UserUsecase,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) domain.AuthUsecase {
	return &authUsecase{
		userUsecase:    userUsecase,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (a *authUsecase) SignUp(ctx context.Context, name string, email string, password string) (res *domain.LoginResponse, err error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))

	user := &domain.User{
		Name:           name,
		Email:          email,
		PasswordDigest: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	err = a.userUsecase.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return a.SignIn(ctx, email, password)
}

func (a *authUsecase) SignIn(ctx context.Context, email string, password string) (res *domain.LoginResponse, err error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userUsecase.GetByEmailPassword(ctx, email, password)
	if err != nil {
		return nil, auth.ErrUserNotFound
	}

	claims := domain.AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed_token, err := token.SignedString(a.signingKey)
	res = &domain.LoginResponse{
		Token: signed_token,
		Id:    user.Id,
	}
	return res, err
}

func (a *authUsecase) ParseToken(ctx context.Context, accessToken string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &domain.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
