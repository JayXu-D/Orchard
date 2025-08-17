package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
)

// SysAlbum 相册表
type SysAlbum struct {
	global.GVA_MODEL
	CreatorUUID   uuid.UUID `json:"creatorUUID" gorm:"index;comment:创建者UUID"`                                                  // 创建者UUID
	Title         string    `json:"title" gorm:"comment:相册标题"`                                                                 // 相册标题
	CoverImageURL string    `json:"coverImageURL" gorm:"comment:相册封面图URL"`                                                     // 相册封面图URL
	Description   string    `json:"description" gorm:"comment:相册描述"`                                                           // 相册描述
	Status        int       `json:"status" gorm:"default:1;comment:相册状态 1:正常 2:禁用"`                                            // 相册状态
	Creator       SysUser   `json:"creator" gorm:"foreignKey:CreatorUUID;references:UUID;comment:创建者信息"`                       // 创建者信息
	AdminUserIDs  []uint    `json:"adminUserIDs" gorm:"-"`                                                                     // 管理员ID列表（用于接收前端数据）
	AdminUsers    []SysUser `json:"adminUsers" gorm:"many2many:sys_album_admin;joinForeignKey:AlbumID;joinReferences:UserID;"` // 管理员列表
}

// SysAlbumAdmin 相册管理员关联表
type SysAlbumAdmin struct {
	AlbumID uint `json:"albumID" gorm:"primaryKey;comment:相册ID"`
	UserID  uint `json:"userID" gorm:"primaryKey;comment:用户ID"`
}

func (SysAlbum) TableName() string {
	return "sys_albums"
}

func (SysAlbumAdmin) TableName() string {
	return "sys_album_admin"
}
