package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
)

// SysDownloadHistory 下载历史记录结构体
type SysDownloadHistory struct {
	global.GVA_MODEL
	UserUUID   uuid.UUID  `json:"userUUID" gorm:"index;comment:用户UUID"`                           // 用户UUID
	DrawingID  uint       `json:"drawingId" gorm:"index;comment:图纸ID"`                            // 图纸ID
	AlbumID    uint       `json:"albumId" gorm:"index;comment:相册ID"`                              // 相册ID
	DownloadAt int64      `json:"downloadAt" gorm:"comment:下载时间戳"`                                // 下载时间戳
	User       SysUser    `json:"user" gorm:"foreignKey:UserUUID;references:UUID;comment:用户信息"`   // 用户信息
	Drawing    SysDrawing `json:"drawing" gorm:"foreignKey:DrawingID;references:ID;comment:图纸信息"` // 图纸信息
	Album      SysAlbum   `json:"album" gorm:"foreignKey:AlbumID;references:ID;comment:相册信息"`     // 相册信息
}

// TableName 下载历史记录表名
func (SysDownloadHistory) TableName() string {
	return "sys_download_histories"
}
