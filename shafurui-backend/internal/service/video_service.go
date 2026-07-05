package service

import (
	"fmt"
	"strings"
	"time"

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

func (s *VideoService) RefreshVideos(videoDirPath string) (*model.VideoRefreshResponse, error) {
	if strings.TrimSpace(videoDirPath) == "" {
		return nil, fmt.Errorf("video_dir_path 未配置")
	}

	start := time.Now()
	result, err := s.repo.RefreshVideos(videoDirPath)
	if err != nil {
		return nil, err
	}

	return &model.VideoRefreshResponse{
		ScannedCount: countVideoItems(result),
		Duration:     time.Since(start).String(),
	}, nil
}

func countVideoItems(result *model.VideoListResponse) int {
	if result == nil {
		return 0
	}

	count := 0
	for _, group := range result.Groups {
		count += len(group.Items)
	}
	return count
}
