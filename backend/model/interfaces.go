package model

import (
	"context"
	"time"
)

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID uint, tokenID uint, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID uint, prevTokenID uint) error
	DeleteUserRefreshTokens(ctx context.Context, userID uint) error
}

type PipelineRepository interface {
	FindByID(id uint) (*Pipeline, error)
	FindByOwner(ownerId uint) ([]Pipeline, error)
	Create(u *Pipeline) error
	Update(u *Pipeline) error
	Delete(id uint) error
}

type RunRepository interface {
	FindByID(id uint) (*Run, error)
	FindByPipeline(pipelineId uint) ([]Run, error)
	Create(u *Run) error
	Update(u *Run) error
	Delete(id uint) error
}
