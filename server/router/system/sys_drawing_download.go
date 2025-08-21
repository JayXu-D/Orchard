package system

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DrawingDownloadRouter struct{}

// InitDrawingDownloadRouter 初始化图纸下载路由
func (s *DrawingDownloadRouter) InitDrawingDownloadRouter(Router *gin.RouterGroup) {
	drawingDownloadRouter := Router.Group("drawing").Use(middleware.OperationRecord())
	{
		drawingDownloadRouter.GET("download/:filename", s.DownloadFile) // 下载文件
	}
}

// DownloadFile 下载文件
func (s *DrawingDownloadRouter) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 构建文件路径
	filePath := filepath.Join("cache", "watermark", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(filePath)
}
