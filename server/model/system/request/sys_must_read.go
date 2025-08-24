package request

import (
	"github.com/google/uuid"
)

// CreateMustRead 创建必读内容请求结构
type CreateMustRead struct {
	CreatorUUID uuid.UUID `json:"creatorUUID" binding:"required" example:"创建者UUID"`
	Title       string    `json:"title" binding:"required" example:"必读内容标题"`
	Content     string    `json:"content" binding:"required" example:"必读内容"`
}

// UpdateMustRead 更新必读内容请求结构
type UpdateMustRead struct {
	ID      uint   `json:"id" binding:"required" example:"必读内容ID"`
	Title   string `json:"title" example:"必读内容标题"`
	Content string `json:"content" example:"必读内容"`
	Status  int    `json:"status" example:"状态"`
}

// GetMustReadByID 根据ID获取必读内容请求结构
type GetMustReadByID struct {
	ID uint `json:"id" binding:"required" example:"必读内容ID"`
}

// DeleteMustRead 删除必读内容请求结构
type DeleteMustRead struct {
	ID uint `json:"id" binding:"required" example:"必读内容ID"`
}

// GetLatestMustRead 获取最新必读内容请求结构
type GetLatestMustRead struct {
	// 无需参数，获取最新的必读内容
}
