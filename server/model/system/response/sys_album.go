package response

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/google/uuid"
)

// AlbumResponse 相册响应结构
type AlbumResponse struct {
	ID            uint       `json:"id" example:"相册ID"`
	CreatorUUID   uuid.UUID  `json:"creatorUUID" example:"创建者UUID"`
	Title         string     `json:"title" example:"相册标题"`
	CoverImageURL string     `json:"coverImageURL" example:"封面图URL"`
	Description   string     `json:"description" example:"相册描述"`
	Status        int        `json:"status" example:"相册状态"`
	CreatedAt     time.Time  `json:"createdAt" example:"创建时间"`
	UpdatedAt     time.Time  `json:"updatedAt" example:"更新时间"`
	Creator       UserInfo   `json:"creator" example:"创建者信息"`
	AdminUsers    []UserInfo `json:"adminUsers" example:"管理员列表"`
}

// UserInfo 用户信息响应结构
type UserInfo struct {
	ID        uint      `json:"id" example:"用户ID"`
	UUID      uuid.UUID `json:"uuid" example:"用户UUID"`
	Username  string    `json:"username" example:"用户名"`
	NickName  string    `json:"nickName" example:"昵称"`
	HeaderImg string    `json:"headerImg" example:"头像"`
}

// AlbumListResponse 相册列表响应结构
type AlbumListResponse struct {
	Albums []AlbumResponse `json:"albums" example:"相册列表"`
	Total  int64           `json:"total" example:"总数"`
}

// 将系统相册模型转换为响应结构
func ToAlbumResponse(album system.SysAlbum) AlbumResponse {
	response := AlbumResponse{
		ID:            album.ID,
		CreatorUUID:   album.CreatorUUID,
		Title:         album.Title,
		CoverImageURL: album.CoverImageURL,
		Description:   album.Description,
		Status:        album.Status,
		CreatedAt:     album.CreatedAt,
		UpdatedAt:     album.UpdatedAt,
	}

	// 转换创建者信息
	if album.Creator.ID != 0 {
		response.Creator = UserInfo{
			ID:        album.Creator.ID,
			UUID:      album.Creator.UUID,
			Username:  album.Creator.Username,
			NickName:  album.Creator.NickName,
			HeaderImg: album.Creator.HeaderImg,
		}
	}

	// 转换管理员信息
	if len(album.AdminUsers) > 0 {
		response.AdminUsers = make([]UserInfo, len(album.AdminUsers))
		for i, admin := range album.AdminUsers {
			response.AdminUsers[i] = UserInfo{
				ID:        admin.ID,
				UUID:      admin.UUID,
				Username:  admin.Username,
				NickName:  admin.NickName,
				HeaderImg: admin.HeaderImg,
			}
		}
	}

	return response
}

// 将系统相册模型列表转换为响应结构
func ToAlbumListResponse(albums []system.SysAlbum, total int64) AlbumListResponse {
	albumResponses := make([]AlbumResponse, len(albums))
	for i, album := range albums {
		albumResponses[i] = ToAlbumResponse(album)
	}

	return AlbumListResponse{
		Albums: albumResponses,
		Total:  total,
	}
}
