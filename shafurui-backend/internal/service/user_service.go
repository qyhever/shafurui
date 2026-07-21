package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"shafurui/internal/domain"
	"shafurui/internal/model"
	"shafurui/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetCurrentUserInfo(ctx context.Context, userID uint64) (*model.UserInfoResponse, error) {
	if userID == 0 {
		return nil, domain.ErrInvalidAdminCredentials
	}

	user, err := s.repo.FindEnabledByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrInvalidAdminCredentials
	}
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrInvalidAdminCredentials
	}

	return &model.UserInfoResponse{
		UserID:   int64(user.ID),
		Username: strings.TrimSpace(user.Username),
		Nickname: strings.TrimSpace(user.Nickname),
	}, nil
}
