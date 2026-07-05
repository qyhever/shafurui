package repository

import "shafurui/internal/model"

type VideoRepository interface {
	ListVideos(videoDirPath string) (*model.VideoListResponse, error)
	RefreshVideos(videoDirPath string) (*model.VideoListResponse, error)
}
