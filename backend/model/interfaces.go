package model

import (
	"context"
	"time"
)

type UserService interface {
	Get(id uint) (*User, error)
	Signup(username string, password string) error
	Signin(user *User) error
	UpdateDetails(user *User) error
}

type UserRepository interface {
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	Create(u *User) error
	Update(u *User) error
}

type TokenService interface {
	NewFirstPairFromUser(ctx context.Context, u *User) (*TokenPair, error)
	NewPairFromUser(ctx context.Context, u *User, prevTokenID uint) (*TokenPair, error)
	Signout(ctx context.Context, uid uint) error
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error)
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID uint, tokenID uint, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID uint, prevTokenID uint) error
	DeleteUserRefreshTokens(ctx context.Context, userID uint) error
}

type PipelineService interface {
	Get(id uint) (*Pipeline, error)
	GetByOwner(ownerId uint) ([]Pipeline, error)
	Create(userId uint, name string, definition string) error
	Update(pipeline *Pipeline) error
	Delete(id uint) error
}

type PipelineRepository interface {
	FindByID(id uint) (*Pipeline, error)
	FindByOwner(ownerId uint) ([]Pipeline, error)
	Create(u *Pipeline) error
	Update(u *Pipeline) error
	Delete(id uint) error
}
