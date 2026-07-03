package service

import (
	"shafurui/internal/model"
	"shafurui/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetCurrentUserInfo() (*model.UserInfoResponse, error) {
	return s.repo.GetCurrentUserInfo()
}
