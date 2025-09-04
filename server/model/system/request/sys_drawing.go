package request

import (
	"github.com/google/uuid"
)

// CreateDrawing 创建图纸请求
type CreateDrawing struct {
	AlbumID            uint      `json:"albumId" binding:"required"`        // 相册ID
	SerialNumber       string    `json:"serialNumber" binding:"required"`   // 图纸序号
	Name               string    `json:"name" binding:"required"`           // 图纸名称
	BeanQuantity       *int      `json:"beanQuantity"`                      // 豆量
	PosterImageURL     string    `json:"posterImageURL" binding:"required"` // 海报图URL
	DrawingURLs        []string  `json:"drawingURLs" binding:"required"`    // 图纸文件URLs
	CreatorUUID        uuid.UUID `json:"creatorUUID" binding:"required"`    // 创建者UUID
	AllowedMemberUUIDs []string  `json:"allowedMemberUUIDs"`                // 允许下载的成员UUIDs
}

// UpdateDrawing 更新图纸请求
type UpdateDrawing struct {
	ID                 uint     `json:"id" binding:"required"`             // 图纸ID
	AlbumID            uint     `json:"albumId" binding:"required"`        // 相册ID
	SerialNumber       string   `json:"serialNumber" binding:"required"`   // 图纸序号
	Name               string   `json:"name" binding:"required"`           // 图纸名称
	BeanQuantity       *int     `json:"beanQuantity"`                      // 豆量
	PosterImageURL     string   `json:"posterImageURL" binding:"required"` // 海报图URL
	DrawingURLs        []string `json:"drawingURLs" binding:"required"`    // 图纸文件URLs
	AllowedMemberUUIDs []string `json:"allowedMemberUUIDs"`                // 允许下载的成员UUIDs
}

// DeleteDrawing 删除图纸请求
type DeleteDrawing struct {
	ID uint `json:"id" binding:"required"` // 图纸ID
}

// GetDrawingByID 根据ID获取图纸请求
type GetDrawingByID struct {
	ID uint `json:"id" binding:"required"` // 图纸ID
}

// GetDrawingList 获取图纸列表请求
type GetDrawingList struct {
	AlbumID   uint   `json:"albumId" binding:"required"` // 相册ID
	Page      int    `json:"page"`                       // 页码
	PageSize  int    `json:"pageSize"`                   // 每页大小
	Keyword   string `json:"keyword"`                    // 搜索关键词
	CreatorID uint   `json:"creatorId"`                  // 创建者ID
}

// GetMyDrawings 获取当前用户可下载的图纸列表请求
type GetMyDrawings struct {
	Page     int    `json:"page"`     // 页码
	PageSize int    `json:"pageSize"` // 每页大小
	Keyword  string `json:"keyword"`  // 搜索关键词
	UserID   uint   `json:"-"`        // 当前用户ID（从JWT中获取）
	UserUUID string `json:"-"`        // 当前用户UUID（从JWT中获取）
}

// DownloadDrawing 下载图纸请求
type DownloadDrawing struct {
	DrawingID     uint   `json:"drawingId" binding:"required"` // 图纸ID
	AlbumID       uint   `json:"albumId" binding:"required"`   // 相册ID
	AddWatermark  bool   `json:"addWatermark"`                 // 是否添加水印
	WatermarkText string `json:"watermarkText"`                // 水印文字
}

// BatchDownloadDrawings 批量下载图纸请求
type BatchDownloadDrawings struct {
	DrawingIDs    []uint `json:"drawingIds" binding:"required"` // 图纸ID列表
	AlbumID       uint   `json:"albumId" binding:"required"`    // 相册ID
	AddWatermark  bool   `json:"addWatermark"`                  // 是否添加水印
	WatermarkText string `json:"watermarkText"`                 // 水印文字
}

// RecordDownload 记录下载请求
type RecordDownload struct {
	DrawingID uint `json:"drawingId" binding:"required"` // 图纸ID
	AlbumID   uint `json:"albumId" binding:"required"`   // 相册ID
}

// DownloadStatusRequest 批量查询下载状态
type DownloadStatusRequest struct {
	DrawingIDs []uint `json:"drawingIds" binding:"required"`
}
