package service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"shafurui/internal/config"
	"shafurui/internal/model"
	passwordpkg "shafurui/internal/pkg/password"
)

type fakeTelegramSender struct {
	mu       sync.Mutex
	messages []string
	sent     chan string
	err      error
}

func (f *fakeTelegramSender) SendMessage(_ context.Context, text string) error {
	f.mu.Lock()
	f.messages = append(f.messages, text)
	f.mu.Unlock()
	if f.sent != nil {
		f.sent <- text
	}
	return f.err
}

func (f *fakeTelegramSender) messageCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.messages)
}

type fakeUserRepository struct {
	usersByLogin map[string]*model.User
	usersByID    map[uint64]*model.User
	err          error
}

func (r *fakeUserRepository) FindEnabledByLogin(_ context.Context, login string) (*model.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	user, ok := r.usersByLogin[login]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (r *fakeUserRepository) FindEnabledByID(_ context.Context, userID uint64) (*model.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	user, ok := r.usersByID[userID]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func setupAuthServiceConfig(t *testing.T) {
	t.Helper()

	previous := config.GlobalConfig
	config.GlobalConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test-secret",
			AccessExpiresIn:  "1h",
			RefreshExpiresIn: "24h",
		},
	}
	t.Cleanup(func() {
		config.GlobalConfig = previous
	})
}

func newAuthTestUser(t *testing.T) *model.User {
	t.Helper()

	hashedPassword, err := passwordpkg.Hash("correct-password")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	return &model.User{
		ID:       123,
		Username: "admin",
		Nickname: "Admin",
		Password: hashedPassword,
		Email:    "admin@example.com",
	}
}

func newAuthTestRepository(user *model.User) *fakeUserRepository {
	return &fakeUserRepository{
		usersByLogin: map[string]*model.User{
			user.Username: user,
			user.Email:    user,
		},
		usersByID: map[uint64]*model.User{
			user.ID: user,
		},
	}
}

func TestAuthLoginWithUsername(t *testing.T) {
	setupAuthServiceConfig(t)

	user := newAuthTestUser(t)
	service := NewAuthService(newAuthTestRepository(user), nil)

	resp, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
		Username: " admin ",
		Password: " correct-password ",
	})
	if err != nil {
		t.Fatalf("AuthLogin() error = %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("AuthLogin() returned empty token")
	}
}

func TestAuthLoginWithEmail(t *testing.T) {
	setupAuthServiceConfig(t)

	user := newAuthTestUser(t)
	service := NewAuthService(newAuthTestRepository(user), nil)

	resp, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
		Username: "admin@example.com",
		Password: "correct-password",
	})
	if err != nil {
		t.Fatalf("AuthLogin() error = %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("AuthLogin() returned empty token")
	}
}

func TestAuthLoginSendsTelegramNotification(t *testing.T) {
	setupAuthServiceConfig(t)

	user := newAuthTestUser(t)
	sender := &fakeTelegramSender{sent: make(chan string, 1)}
	service := NewAuthService(newAuthTestRepository(user), sender)

	resp, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
		Username: "admin",
		Password: "correct-password",
	})
	if err != nil {
		t.Fatalf("AuthLogin() error = %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("AuthLogin() returned empty token")
	}
	var message string
	select {
	case message = <-sender.sent:
	case <-time.After(time.Second):
		t.Fatal("telegram message was not sent")
	}

	for _, want := range []string{
		"sfr: 用户登录成功",
		"用户名: admin",
		"用户ID: 123",
		"时间: ",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("telegram message = %q, want substring %q", message, want)
		}
	}
	if !regexp.MustCompile(`时间: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`).MatchString(message) {
		t.Fatalf("telegram message time has unexpected format: %q", message)
	}
	if strings.Contains(message, "correct-password") ||
		strings.Contains(message, resp.AccessToken) ||
		strings.Contains(message, resp.RefreshToken) {
		t.Fatalf("telegram message exposed sensitive data: %q", message)
	}
}

func TestAuthLoginIgnoresTelegramNotificationError(t *testing.T) {
	setupAuthServiceConfig(t)

	user := newAuthTestUser(t)
	sender := &fakeTelegramSender{sent: make(chan string, 1), err: errors.New("send failed")}
	service := NewAuthService(newAuthTestRepository(user), sender)

	resp, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
		Username: "admin",
		Password: "correct-password",
	})
	if err != nil {
		t.Fatalf("AuthLogin() error = %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("AuthLogin() returned empty token")
	}
	select {
	case <-sender.sent:
	case <-time.After(time.Second):
		t.Fatal("telegram message was not sent")
	}
}

func TestAuthLoginDoesNotWaitForTelegramNotification(t *testing.T) {
	setupAuthServiceConfig(t)

	user := newAuthTestUser(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	sender := telegramMessageSenderFunc(func(context.Context, string) error {
		close(entered)
		<-release
		return nil
	})
	t.Cleanup(func() {
		close(release)
	})

	service := NewAuthService(newAuthTestRepository(user), sender)
	done := make(chan struct{})
	go func() {
		resp, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
			Username: "admin",
			Password: "correct-password",
		})
		if err != nil {
			t.Errorf("AuthLogin() error = %v", err)
		}
		if resp == nil || resp.AccessToken == "" || resp.RefreshToken == "" {
			t.Error("AuthLogin() returned empty token")
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("AuthLogin() waited for telegram notification")
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("telegram notification did not start")
	}
}

func TestAuthLoginDoesNotSendTelegramNotificationForInvalidCredentials(t *testing.T) {
	setupAuthServiceConfig(t)

	user := newAuthTestUser(t)
	sender := &fakeTelegramSender{}
	service := NewAuthService(newAuthTestRepository(user), sender)

	_, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
		Username: "admin",
		Password: "wrong-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("AuthLogin() error = %v, want %v", err, ErrInvalidCredentials)
	}
	if count := sender.messageCount(); count != 0 {
		t.Fatalf("telegram messages = %d, want 0", count)
	}
}

func TestAuthLoginReturnsInvalidCredentialsForMissingUser(t *testing.T) {
	setupAuthServiceConfig(t)

	sender := &fakeTelegramSender{}
	service := NewAuthService(&fakeUserRepository{
		usersByLogin: map[string]*model.User{},
		usersByID:    map[uint64]*model.User{},
	}, sender)

	_, err := service.AuthLogin(context.Background(), model.AuthLoginRequest{
		Username: "disabled@example.com",
		Password: "correct-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("AuthLogin() error = %v, want %v", err, ErrInvalidCredentials)
	}
	if count := sender.messageCount(); count != 0 {
		t.Fatalf("telegram messages = %d, want 0", count)
	}
}

type telegramMessageSenderFunc func(context.Context, string) error

func (f telegramMessageSenderFunc) SendMessage(ctx context.Context, text string) error {
	return f(ctx, text)
}
