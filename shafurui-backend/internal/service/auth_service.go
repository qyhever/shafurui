package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"shafurui/internal/config"
	"shafurui/internal/domain"
	"shafurui/internal/model"
	jwtpkg "shafurui/internal/pkg/jwt"

	"go.uber.org/zap"
)

type TelegramMessageSender interface {
	SendMessage(ctx context.Context, text string) error
}

type AuthService struct {
	telegramSender TelegramMessageSender
}

func NewAuthService(telegramSender TelegramMessageSender) *AuthService {
	return &AuthService{
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

	cfg := config.GetConfig()
	if cfg == nil {
		return nil, ErrInvalidCredentials
	}

	user := cfg.Auth.DefaultUser
	if user.UserID <= 0 ||
		strings.TrimSpace(user.Username) == "" ||
		user.Password == "" ||
		username != strings.TrimSpace(user.Username) ||
		password != user.Password {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := jwtpkg.GenToken(uint64(user.UserID))
	if err != nil {
		return nil, err
	}
	s.notifyLoginSuccessAsync(ctx, user)
	return &model.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) notifyLoginSuccessAsync(ctx context.Context, user config.DefaultUserConfig) {
	if s.telegramSender == nil {
		return
	}

	message := strings.Join([]string{
		"sfr: 用户登录成功",
		"用户名: " + strings.TrimSpace(user.Username),
		"用户ID: " + strconv.FormatInt(user.UserID, 10),
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
