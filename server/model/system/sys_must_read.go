package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
)

// SysMustRead 必读内容表
type SysMustRead struct {
	global.GVA_MODEL
	CreatorUUID uuid.UUID `json:"creatorUUID" gorm:"index;comment:创建者UUID"`                            // 创建者UUID
	Title       string    `json:"title" gorm:"comment:必读内容标题"`                                         // 必读内容标题
	Content     string    `json:"content" gorm:"type:text;comment:必读内容"`                               // 必读内容
	Status      int       `json:"status" gorm:"default:1;comment:状态 1:正常 2:禁用"`                        // 状态
	Creator     SysUser   `json:"creator" gorm:"foreignKey:CreatorUUID;references:UUID;comment:创建者信息"` // 创建者信息
}

func (SysMustRead) TableName() string {
	return "sys_must_reads"
}
