package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"shafurui/internal/model"
)

func TestUserServiceGetCurrentUserInfo(t *testing.T) {
	user := &model.User{
		ID:       123,
		Username: " admin ",
		Nickname: " Admin ",
	}
	service := NewUserService(&fakeUserRepository{
		usersByLogin: map[string]*model.User{},
		usersByID: map[uint64]*model.User{
			user.ID: user,
		},
	})

	got, err := service.GetCurrentUserInfo(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetCurrentUserInfo() error = %v", err)
	}
	if got.UserID != 123 || got.Username != "admin" || got.Nickname != "Admin" {
		t.Fatalf("GetCurrentUserInfo() = %+v, want trimmed database user", got)
	}
}

func TestUserServiceGetCurrentUserInfoReturnsErrorForMissingUser(t *testing.T) {
	service := NewUserService(&fakeUserRepository{
		usersByLogin: map[string]*model.User{},
		usersByID:    map[uint64]*model.User{},
	})

	_, err := service.GetCurrentUserInfo(context.Background(), 123)
	if err == nil {
		t.Fatal("GetCurrentUserInfo() error = nil, want error")
	}
}

func TestUserServiceGetCurrentUserInfoPropagatesRepositoryError(t *testing.T) {
	wantErr := errors.New("database failed")
	service := NewUserService(&fakeUserRepository{err: wantErr})

	_, err := service.GetCurrentUserInfo(context.Background(), 123)
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetCurrentUserInfo() error = %v, want %v", err, wantErr)
	}
}

func TestFakeUserRepositoryMissingIDUsesSQLNoRows(t *testing.T) {
	repo := &fakeUserRepository{
		usersByLogin: map[string]*model.User{},
		usersByID:    map[uint64]*model.User{},
	}

	_, err := repo.FindEnabledByID(context.Background(), 123)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("FindEnabledByID() error = %v, want %v", err, sql.ErrNoRows)
	}
}
