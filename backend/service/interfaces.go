package service

import (
	"context"
	"di/model"
)

type Services struct {
	UserService     UserService
	PipelineService PipelineService
	TokenService    TokenService
}

type UserService interface {
	Get(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Signup(username string, password string) error
	Signin(user *model.User) error
	UpdateDetails(user *model.User) error
}

type TokenService interface {
	NewFirstPairFromUser(ctx context.Context, u *model.User) (*model.TokenPair, error)
	NewPairFromUser(ctx context.Context, u *model.User, prevTokenID uint) (*model.TokenPair, error)
	Signout(ctx context.Context, uid uint) error
	ValidateIDToken(tokenString string) (*model.User, error)
	ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error)
}

type PipelineService interface {
	Get(id uint) (*model.Pipeline, error)
	GetByOwner(ownerId uint) ([]model.Pipeline, error)
	Create(userId uint, name string, definition string) error
	Update(pipeline *model.Pipeline) error
	Delete(id uint) error
}