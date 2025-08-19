package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
)

// SysDrawing 图纸结构体
type SysDrawing struct {
	global.GVA_MODEL
	AlbumID        uint      `json:"albumId" gorm:"index;comment:相册ID"`                                   // 相册ID
	SerialNumber   string    `json:"serialNumber" gorm:"index;comment:图纸序号"`                              // 图纸序号
	Name           string    `json:"name" gorm:"comment:图纸名称"`                                            // 图纸名称
	BeanQuantity   *int      `json:"beanQuantity" gorm:"comment:豆量"`                                      // 豆量
	PosterImageURL string    `json:"posterImageURL" gorm:"comment:海报图URL"`                                // 海报图URL
	DrawingURLs    string    `json:"drawingURLs" gorm:"type:text;comment:图纸文件URLs"`                       // 图纸文件URLs (JSON格式)
	CreatorUUID    uuid.UUID `json:"creatorUUID" gorm:"index;comment:创建者UUID"`                            // 创建者UUID
	AllowedMembers string    `json:"allowedMembers" gorm:"type:text;comment:允许下载的成员"`                     // 允许下载的成员 (JSON格式)
	Album          SysAlbum  `json:"album" gorm:"foreignKey:AlbumID;references:ID;comment:相册信息"`          // 相册信息
	Creator        SysUser   `json:"creator" gorm:"foreignKey:CreatorUUID;references:UUID;comment:创建者信息"` // 创建者信息
}

// TableName 图纸表名
func (SysDrawing) TableName() string {
	return "sys_drawings"
}
