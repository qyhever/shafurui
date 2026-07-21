package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"shafurui/internal/domain"
	"shafurui/internal/model"
	jwtpkg "shafurui/internal/pkg/jwt"
	passwordpkg "shafurui/internal/pkg/password"
	"shafurui/internal/repository"

	"go.uber.org/zap"
)

type TelegramMessageSender interface {
	SendMessage(ctx context.Context, text string) error
}

type AuthService struct {
	userRepo       repository.UserRepository
	telegramSender TelegramMessageSender
}

func NewAuthService(userRepo repository.UserRepository, telegramSender TelegramMessageSender) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		telegramSender: telegramSender,
	}
}

var (
	ErrInvalidCredentials  = domain.ErrInvalidAdminCredentials
	ErrInvalidRefreshToken = domain.ErrInvalidRefreshToken
)

func (s *AuthService) AuthLogin(ctx context.Context, req model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindEnabledByLogin(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 || strings.TrimSpace(user.Password) == "" {
		return nil, ErrInvalidCredentials
	}
	if err := passwordpkg.Compare(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := jwtpkg.GenToken(user.ID)
	if err != nil {
		return nil, err
	}

	if user.Username == "charon" {
		s.notifyLoginSuccessAsync(ctx, user)
	}
	return &model.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) notifyLoginSuccessAsync(ctx context.Context, user *model.User) {
	if s.telegramSender == nil {
		return
	}

	message := strings.Join([]string{
		"sfr: 用户登录成功",
		"用户名: " + strings.TrimSpace(user.Username),
		"用户ID: " + strconv.FormatUint(user.ID, 10),
		"时间: " + time.Now().Format("2006-01-02 15:04:05"),
	}, "\n")
	go func() {
		if err := s.telegramSender.SendMessage(context.WithoutCancel(ctx), message); err != nil {
			zap.L().Warn("send admin login telegram notification failed", zap.Error(err))
		}
	}()
}

func (s *AuthService) AuthRefreshToken(ctx context.Context, req model.AuthRefreshTokenRequest) (*model.AuthLoginResponse, error) {
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	claims, err := jwtpkg.ParseToken(refreshToken)
	if err != nil || !claims.IsRefreshToken() {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, newRefreshToken, err := jwtpkg.GenToken(claims.UserID)
	if err != nil {
		return nil, err
	}
	return &model.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
