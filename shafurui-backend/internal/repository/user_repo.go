package repository

import "shafurui/internal/model"

type UserRepository interface {
	GetCurrentUserInfo() (*model.UserInfoResponse, error)
}
