package response

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/google/uuid"
)

// MustReadResponse 必读内容响应结构
type MustReadResponse struct {
	ID          uint      `json:"id" example:"必读内容ID"`
	CreatorUUID uuid.UUID `json:"creatorUUID" example:"创建者UUID"`
	Title       string    `json:"title" example:"必读内容标题"`
	Content     string    `json:"content" example:"必读内容"`
	Status      int       `json:"status" example:"状态"`
	CreatedAt   time.Time `json:"createdAt" example:"创建时间"`
	UpdatedAt   time.Time `json:"updatedAt" example:"更新时间"`
	Creator     UserInfo  `json:"creator" example:"创建者信息"`
}

// MustReadListResponse 必读内容列表响应结构
type MustReadListResponse struct {
	MustReads []MustReadResponse `json:"mustReads" example:"必读内容列表"`
	Total     int64              `json:"total" example:"总数"`
}

// 将系统必读内容模型转换为响应结构
func ToMustReadResponse(mustRead system.SysMustRead) MustReadResponse {
	response := MustReadResponse{
		ID:          mustRead.ID,
		CreatorUUID: mustRead.CreatorUUID,
		Title:       mustRead.Title,
		Content:     mustRead.Content,
		Status:      mustRead.Status,
		CreatedAt:   mustRead.CreatedAt,
		UpdatedAt:   mustRead.UpdatedAt,
	}

	// 转换创建者信息
	if mustRead.Creator.ID != 0 {
		response.Creator = UserInfo{
			ID:        mustRead.Creator.ID,
			UUID:      mustRead.Creator.UUID,
			Username:  mustRead.Creator.Username,
			NickName:  mustRead.Creator.NickName,
			HeaderImg: mustRead.Creator.HeaderImg,
		}
	}

	return response
}

// 将系统必读内容模型列表转换为响应结构
func ToMustReadListResponse(mustReads []system.SysMustRead, total int64) MustReadListResponse {
	mustReadResponses := make([]MustReadResponse, len(mustReads))
	for i, mustRead := range mustReads {
		mustReadResponses[i] = ToMustReadResponse(mustRead)
	}

	return MustReadListResponse{
		MustReads: mustReadResponses,
		Total:     total,
	}
}
