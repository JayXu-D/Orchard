package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

var Album = new(album)

type album struct{}

func (a *album) Init() error {
	// 这里可以添加一些初始化的相册数据
	// 例如：创建一些示例相册

	// 检查是否已经有相册数据
	var count int64
	global.GVA_DB.Model(&system.SysAlbum{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 创建示例相册数据（可选）
	// 注意：这里需要确保有对应的用户存在
	/*
		albums := []system.SysAlbum{
			{
				CreatorUUID:   uuid.MustParse("your-user-uuid-here"), // 需要替换为实际的用户UUID
				Title:         "示例相册1",
				CoverImageURL: "https://example.com/cover1.jpg",
				Description:   "这是一个示例相册",
				Status:        1,
			},
			{
				CreatorUUID:   uuid.MustParse("your-user-uuid-here"), // 需要替换为实际的用户UUID
				Title:         "示例相册2",
				CoverImageURL: "https://example.com/cover2.jpg",
				Description:   "这是另一个示例相册",
				Status:        1,
			},
		}

		for _, album := range albums {
			if err := global.GVA_DB.Create(&album).Error; err != nil {
				return err
			}
		}
	*/

	return nil
}
