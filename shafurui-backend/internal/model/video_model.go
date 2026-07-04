package model

import "time"

type VideoListResponse struct {
	Groups []VideoGroup `json:"groups"`
}

type VideoGroup struct {
	Date  string      `json:"date"`
	Items []VideoItem `json:"items"`
}

type VideoItem struct {
	ID           string    `json:"id"`
	Filename     string    `json:"filename"`
	RelativePath string    `json:"relativePath"`
	URL          string    `json:"url"`
	CoverURL     string    `json:"coverUrl"`
	ShotAt       time.Time `json:"shotAt"`
	GroupDate    string    `json:"groupDate"`
	DurationSec  *float64  `json:"durationSec,omitempty"`
	Width        *int      `json:"width,omitempty"`
	Height       *int      `json:"height,omitempty"`
	SizeBytes    int64     `json:"sizeBytes"`
	Mtime        time.Time `json:"mtime"`
}
