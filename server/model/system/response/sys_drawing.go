package response

import (
	"encoding/json"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/google/uuid"
)

// DrawingResponse 图纸响应结构体
type DrawingResponse struct {
	ID                 uint      `json:"id"`                 // 图纸ID
	AlbumID            uint      `json:"albumId"`            // 相册ID
	SerialNumber       string    `json:"serialNumber"`       // 图纸序号
	Name               string    `json:"name"`               // 图纸名称
	BeanQuantity       *int      `json:"beanQuantity"`       // 豆量
	PosterImageURL     string    `json:"posterImageURL"`     // 海报图URL
	DrawingURLs        []string  `json:"drawingURLs"`        // 图纸文件URLs
	CreatorUUID        uuid.UUID `json:"creatorUUID"`        // 创建者UUID
	AllowedMemberUUIDs []string  `json:"allowedMemberUUIDs"` // 允许下载的成员UUIDs
	CreatedAt          string    `json:"createdAt"`          // 创建时间
	UpdatedAt          string    `json:"updatedAt"`          // 更新时间
	Album              struct {
		ID    uint   `json:"id"`    // 相册ID
		Title string `json:"title"` // 相册标题
	} `json:"album"` // 相册信息
	Creator struct {
		UUID     uuid.UUID `json:"uuid"`     // 用户UUID
		Username string    `json:"username"` // 用户名
		NickName string    `json:"nickName"` // 昵称
	} `json:"creator"` // 创建者信息
}

// DrawingListResponse 图纸列表响应结构体
type DrawingListResponse struct {
	Drawings []DrawingResponse `json:"drawings"` // 图纸列表
	Total    int64             `json:"total"`    // 总数
}

// DownloadResponse 下载响应结构体
type DownloadResponse struct {
	DownloadURL string   `json:"downloadUrl"` // 下载链接
	FileName    string   `json:"fileName"`    // 文件名
	FileSize    int64    `json:"fileSize"`    // 文件大小
	FilePaths   []string `json:"filePaths"`   // 文件路径列表（用于批量下载）
}

// ToDrawingResponse 转换为图纸响应结构体
func ToDrawingResponse(drawing *system.SysDrawing) DrawingResponse {
	var drawingURLs []string
	var allowedMemberUUIDs []string

	// 解析图纸文件URLs
	if drawing.DrawingURLs != "" {
		_ = json.Unmarshal([]byte(drawing.DrawingURLs), &drawingURLs)
	}

	// 解析允许下载的成员UUIDs
	if drawing.AllowedMembers != "" {
		_ = json.Unmarshal([]byte(drawing.AllowedMembers), &allowedMemberUUIDs)
	}

	response := DrawingResponse{
		ID:                 drawing.ID,
		AlbumID:            drawing.AlbumID,
		SerialNumber:       drawing.SerialNumber,
		Name:               drawing.Name,
		BeanQuantity:       drawing.BeanQuantity,
		PosterImageURL:     drawing.PosterImageURL,
		DrawingURLs:        drawingURLs,
		CreatorUUID:        drawing.CreatorUUID,
		AllowedMemberUUIDs: allowedMemberUUIDs,
		CreatedAt:          drawing.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          drawing.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 设置相册信息
	if drawing.Album.ID != 0 {
		response.Album.ID = drawing.Album.ID
		response.Album.Title = drawing.Album.Title
	}

	// 设置创建者信息
	if drawing.Creator.UUID != uuid.Nil {
		response.Creator.UUID = drawing.Creator.UUID
		response.Creator.Username = drawing.Creator.Username
		response.Creator.NickName = drawing.Creator.NickName
	}

	return response
}

// ToDrawingListResponse 转换为图纸列表响应结构体
func ToDrawingListResponse(drawings []*system.SysDrawing, total int64) DrawingListResponse {
	drawingResponses := make([]DrawingResponse, len(drawings))
	for i, drawing := range drawings {
		drawingResponses[i] = ToDrawingResponse(drawing)
	}

	return DrawingListResponse{
		Drawings: drawingResponses,
		Total:    total,
	}
}
