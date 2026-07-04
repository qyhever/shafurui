package service

import (
	"fmt"
	"strings"

	"shafurui/internal/model"
	"shafurui/internal/repository"
)

type VideoService struct {
	repo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) ListVideos(videoDirPath string) (*model.VideoListResponse, error) {
	if strings.TrimSpace(videoDirPath) == "" {
		return nil, fmt.Errorf("video_dir_path 未配置")
	}
	return s.repo.ListVideos(videoDirPath)
}
