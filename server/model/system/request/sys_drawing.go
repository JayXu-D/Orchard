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
