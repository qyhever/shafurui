package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"shafurui/internal/model"
)

func TestVideoRepositoryListVideosReturnsEmptyWhenIndexMissing(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "2024-05-06_070809.mp4"), []byte("video"), 0o644); err != nil {
		t.Fatalf("write video file: %v", err)
	}

	repo := &VideoRepositoryImpl{}
	repo.loadIndex(root)

	result, err := repo.ListVideos(root)
	if err != nil {
		t.Fatalf("ListVideos() error = %v", err)
	}
	if len(result.Groups) != 0 {
		t.Fatalf("ListVideos() groups length = %d, want 0", len(result.Groups))
	}
}

func TestVideoRepositoryRefreshVideosScansAndWritesIndex(t *testing.T) {
	root := t.TempDir()
	videoPath := filepath.Join(root, "2024-05-06_070809.mp4")
	if err := os.WriteFile(videoPath, []byte("video"), 0o644); err != nil {
		t.Fatalf("write video file: %v", err)
	}

	repo := &VideoRepositoryImpl{}
	result, err := repo.RefreshVideos(root)
	if err != nil {
		t.Fatalf("RefreshVideos() error = %v", err)
	}

	assertVideoCount(t, result, 1)
	indexPath := filepath.Join(root, videoIndexFilename)
	if _, err := os.Stat(indexPath); err != nil {
		t.Fatalf("stat index file: %v", err)
	}

	var indexed model.VideoListResponse
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index file: %v", err)
	}
	if err := json.Unmarshal(data, &indexed); err != nil {
		t.Fatalf("unmarshal index file: %v", err)
	}
	assertVideoCount(t, &indexed, 1)
}

func TestVideoRepositoryListVideosUsesCacheAfterRefresh(t *testing.T) {
	root := t.TempDir()
	videoPath := filepath.Join(root, "2024-05-06_070809.mp4")
	if err := os.WriteFile(videoPath, []byte("video"), 0o644); err != nil {
		t.Fatalf("write video file: %v", err)
	}

	repo := &VideoRepositoryImpl{}
	if _, err := repo.RefreshVideos(root); err != nil {
		t.Fatalf("RefreshVideos() error = %v", err)
	}
	if err := os.Remove(videoPath); err != nil {
		t.Fatalf("remove video file: %v", err)
	}

	result, err := repo.ListVideos(root)
	if err != nil {
		t.Fatalf("ListVideos() error = %v", err)
	}
	assertVideoCount(t, result, 1)
}

func TestVideoRepositoryLoadIndexOnStartup(t *testing.T) {
	root := t.TempDir()
	repo := &VideoRepositoryImpl{}
	if err := os.WriteFile(filepath.Join(root, "2024-05-06_070809.mp4"), []byte("video"), 0o644); err != nil {
		t.Fatalf("write video file: %v", err)
	}
	if _, err := repo.RefreshVideos(root); err != nil {
		t.Fatalf("RefreshVideos() error = %v", err)
	}

	reloaded := &VideoRepositoryImpl{}
	reloaded.loadIndex(root)
	result, err := reloaded.ListVideos(root)
	if err != nil {
		t.Fatalf("ListVideos() error = %v", err)
	}
	assertVideoCount(t, result, 1)
}

func assertVideoCount(t *testing.T, result *model.VideoListResponse, want int) {
	t.Helper()

	got := 0
	for _, group := range result.Groups {
		got += len(group.Items)
	}
	if got != want {
		t.Fatalf("video count = %d, want %d", got, want)
	}
}
