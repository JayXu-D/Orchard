package request

import (
	common "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/google/uuid"
)

// CreateAlbum 创建相册请求结构
type CreateAlbum struct {
	CreatorUUID   uuid.UUID `json:"creatorUUID" binding:"required" example:"创建者UUID"`
	Title         string    `json:"title" binding:"required" example:"相册标题"`
	CoverImageURL string    `json:"coverImageURL" example:"封面图URL"`
	Description   string    `json:"description" example:"相册描述"`
	AdminUserIDs  []uint    `json:"adminUserIDs" example:"管理员ID列表"`
}

// UpdateAlbum 更新相册请求结构
type UpdateAlbum struct {
	ID            uint   `json:"id" binding:"required" example:"相册ID"`
	Title         string `json:"title" example:"相册标题"`
	CoverImageURL string `json:"coverImageURL" example:"封面图URL"`
	Description   string `json:"description" example:"相册描述"`
	Status        int    `json:"status" example:"相册状态"`
	AdminUserIDs  []uint `json:"adminUserIDs" example:"管理员ID列表"`
}

// GetAlbumList 获取相册列表请求结构
type GetAlbumList struct {
	common.PageInfo
	Title       string    `json:"title" form:"title" example:"相册标题"`
	CreatorUUID uuid.UUID `json:"creatorUUID" form:"creatorUUID" example:"创建者UUID"`
	Status      int       `json:"status" form:"status" example:"相册状态"`
}

// GetAlbumByID 根据ID获取相册请求结构
type GetAlbumByID struct {
	ID uint `json:"id" binding:"required" example:"相册ID"`
}

// DeleteAlbum 删除相册请求结构
type DeleteAlbum struct {
	ID uint `json:"id" binding:"required" example:"相册ID"`
}
