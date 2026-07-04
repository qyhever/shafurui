package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"shafurui/internal/config"
	"shafurui/internal/model"
	"shafurui/internal/repository"
)

var (
	videoExtensions = map[string]struct{}{
		".mp4":  {},
		".mov":  {},
		".m4v":  {},
		".webm": {},
	}
	pathTimeRegexp = regexp.MustCompile(`(\d{4})[-_/]?(\d{2})[-_/]?(\d{2})(?:[ T_-]?(\d{2})[-_:.]?(\d{2})[-_:.]?(\d{2}))?`)
)

type VideoRepositoryImpl struct{}

func NewVideoRepository() repository.VideoRepository {
	return &VideoRepositoryImpl{}
}

func (r *VideoRepositoryImpl) ListVideos(videoDirPath string) (*model.VideoListResponse, error) {
	root := filepath.Clean(videoDirPath)
	rootInfo, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("读取视频目录失败: %w", err)
	}
	if !rootInfo.IsDir() {
		return nil, fmt.Errorf("视频路径不是目录: %s", videoDirPath)
	}

	var items []model.VideoItem
	if err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isVideoFile(path) {
			return nil
		}

		item, err := buildVideoItem(root, path)
		if err != nil {
			return err
		}
		items = append(items, item)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("扫描视频目录失败: %w", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ShotAt.After(items[j].ShotAt)
	})

	return &model.VideoListResponse{Groups: groupVideoItems(items)}, nil
}

func buildVideoItem(root, path string) (model.VideoItem, error) {
	info, err := os.Stat(path)
	if err != nil {
		return model.VideoItem{}, err
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return model.VideoItem{}, err
	}
	rel = filepath.ToSlash(rel)

	meta := probeVideoMetadata(path)
	shotAt := firstNonZeroTime(meta.CreationTime, parseShotAtFromPath(rel), info.ModTime())
	groupDate := shotAt.Format(time.DateOnly)

	coverRel := replaceExt(rel, ".jpg")
	videoBaseURL := strings.TrimRight(config.GetVideoBaseURL(), "/")
	return model.VideoItem{
		ID:           rel,
		Filename:     filepath.Base(path),
		RelativePath: rel,
		URL:          joinVideoURL(videoBaseURL, rel),
		CoverURL:     joinVideoURL(videoBaseURL, coverRel),
		ShotAt:       shotAt,
		GroupDate:    groupDate,
		DurationSec:  meta.DurationSec,
		Width:        meta.Width,
		Height:       meta.Height,
		SizeBytes:    info.Size(),
		Mtime:        info.ModTime(),
	}, nil
}

func joinVideoURL(baseURL, rel string) string {
	if baseURL == "" {
		return rel
	}
	return baseURL + "/" + strings.TrimLeft(rel, "/")
}

func groupVideoItems(items []model.VideoItem) []model.VideoGroup {
	groups := make([]model.VideoGroup, 0)
	groupIndex := make(map[string]int)
	for _, item := range items {
		idx, ok := groupIndex[item.GroupDate]
		if !ok {
			groupIndex[item.GroupDate] = len(groups)
			groups = append(groups, model.VideoGroup{
				Date:  item.GroupDate,
				Items: []model.VideoItem{},
			})
			idx = len(groups) - 1
		}
		groups[idx].Items = append(groups[idx].Items, item)
	}
	return groups
}

func isVideoFile(path string) bool {
	_, ok := videoExtensions[strings.ToLower(filepath.Ext(path))]
	return ok
}

func replaceExt(path, ext string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + ext
}

func firstNonZeroTime(values ...time.Time) time.Time {
	for _, value := range values {
		if !value.IsZero() {
			return value
		}
	}
	return time.Now()
}

func parseShotAtFromPath(path string) time.Time {
	matches := pathTimeRegexp.FindStringSubmatch(path)
	if len(matches) == 0 {
		return time.Time{}
	}

	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	hour, minute, second := 0, 0, 0
	if matches[4] != "" {
		hour, _ = strconv.Atoi(matches[4])
		minute, _ = strconv.Atoi(matches[5])
		second, _ = strconv.Atoi(matches[6])
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

type videoMetadata struct {
	CreationTime time.Time
	DurationSec  *float64
	Width        *int
	Height       *int
}

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
	Format struct {
		Duration string            `json:"duration"`
		Tags     map[string]string `json:"tags"`
	} `json:"format"`
}

func probeVideoMetadata(path string) videoMetadata {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return videoMetadata{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration:format_tags=creation_time:stream=width,height",
		"-of", "json",
		path,
	).Output()
	if err != nil || len(out) == 0 {
		return videoMetadata{}
	}

	var parsed ffprobeOutput
	if err := json.Unmarshal(out, &parsed); err != nil {
		return videoMetadata{}
	}

	meta := videoMetadata{}
	if duration, err := strconv.ParseFloat(parsed.Format.Duration, 64); err == nil {
		meta.DurationSec = &duration
	}
	if creationTime := parseCreationTime(parsed.Format.Tags); !creationTime.IsZero() {
		meta.CreationTime = creationTime
	}
	for _, stream := range parsed.Streams {
		if stream.Width > 0 && stream.Height > 0 {
			width := stream.Width
			height := stream.Height
			meta.Width = &width
			meta.Height = &height
			break
		}
	}
	return meta
}

func parseCreationTime(tags map[string]string) time.Time {
	if len(tags) == 0 {
		return time.Time{}
	}
	value, ok := tags["creation_time"]
	if !ok {
		return time.Time{}
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006:01:02 15:04:05",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.Local()
		}
	}

	if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local); err == nil {
		return parsed
	}
	if parsed, err := time.ParseInLocation("2006:01:02 15:04:05", value, time.Local); err == nil {
		return parsed
	}
	return time.Time{}
}
