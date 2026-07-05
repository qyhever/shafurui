package controller

import (
	"shafurui/internal/config"
	"shafurui/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoController struct {
	videoService *service.VideoService
}

func NewVideoController(videoService *service.VideoService) *VideoController {
	return &VideoController{
		videoService: videoService,
	}
}

// ListVideos godoc
// @Summary 获取视频列表
// @Description 读取视频索引缓存，按日期分组返回视频列表。
// @Tags video
// @Produce json
// @Success 200 {object} ResponseData
// @Router /video [get]
func (vc *VideoController) ListVideos(c *gin.Context) {
	result, err := vc.videoService.ListVideos(config.GetVideoDirPath())
	if err != nil {
		zap.L().Error("list videos failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, result)
}

// RefreshVideos godoc
// @Summary 刷新视频索引
// @Description 扫描配置的视频目录，生成视频索引文件并更新内存缓存。
// @Tags video
// @Produce json
// @Success 200 {object} ResponseData
// @Router /video/refresh [post]
func (vc *VideoController) RefreshVideos(c *gin.Context) {
	result, err := vc.videoService.RefreshVideos(config.GetVideoDirPath())
	if err != nil {
		zap.L().Error("refresh videos failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, result)
}
