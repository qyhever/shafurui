package repository

import (
	"context"

	"shafurui/internal/model"
)

type UserRepository interface {
	FindEnabledByLogin(ctx context.Context, login string) (*model.User, error)
	FindEnabledByID(ctx context.Context, userID uint64) (*model.User, error)
}
