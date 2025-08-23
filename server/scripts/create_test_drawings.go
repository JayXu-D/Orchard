package main

import (
	"fmt"
	"log"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func main() {
	// 初始化数据库连接
	initialize.Gorm()

	fmt.Println("开始更新测试图纸数据...")

	// 更新现有的空白图纸记录
	var drawings []system.SysDrawing
	err := global.GVA_DB.Find(&drawings).Error
	if err != nil {
		log.Fatalf("查询图纸失败: %v", err)
	}

	fmt.Printf("找到 %d 个图纸记录\n", len(drawings))

	// 检查是否有相册和用户数据
	var albums []system.SysAlbum
	err = global.GVA_DB.Find(&albums).Error
	if err != nil {
		log.Fatalf("查询相册失败: %v", err)
	}

	var users []system.SysUser
	err = global.GVA_DB.Find(&users).Error
	if err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}

	if len(albums) == 0 {
		fmt.Println("没有找到相册数据，创建测试相册...")
		// 创建测试相册
		testAlbum := system.SysAlbum{
			Title:         "测试相册",
			CoverImageURL: "uploads/test/cover.jpg",
			Description:   "这是一个测试相册",
			Status:        1,
		}

		if len(users) > 0 {
			testAlbum.CreatorUUID = users[0].UUID
		}

		err = global.GVA_DB.Create(&testAlbum).Error
		if err != nil {
			log.Fatalf("创建测试相册失败: %v", err)
		}
		albums = append(albums, testAlbum)
		fmt.Printf("创建测试相册成功，ID: %d\n", testAlbum.ID)
	}

	if len(users) == 0 {
		fmt.Println("没有找到用户数据，无法继续...")
		return
	}

	// 使用第一个相册和第一个用户
	album := albums[0]
	user := users[0]

	fmt.Printf("使用相册: %s (ID: %d)\n", album.Title, album.ID)
	fmt.Printf("使用用户: %s (UUID: %s)\n", user.Username, user.UUID)

	// 更新现有的图纸记录
	for i, drawing := range drawings {
		fmt.Printf("更新图纸 %d...\n", drawing.ID)

		// 更新图纸数据
		updates := map[string]interface{}{
			"album_id":         album.ID,
			"serial_number":    fmt.Sprintf("TEST-%03d", i+1),
			"name":             fmt.Sprintf("测试图纸%d", i+1),
			"bean_quantity":    (i + 1) * 100,
			"poster_image_url": fmt.Sprintf("uploads/test/poster%d.jpg", i+1),
			"drawing_urls":     fmt.Sprintf(`["uploads/test/drawing%d.pdf", "uploads/test/drawing%d.dwg"]`, i+1, i+1),
			"creator_uuid":     user.UUID,
			"allowed_members":  fmt.Sprintf(`["%s"]`, user.UUID.String()),
		}

		err = global.GVA_DB.Model(&drawing).Updates(updates).Error
		if err != nil {
			log.Printf("更新图纸 %d 失败: %v", drawing.ID, err)
		} else {
			fmt.Printf("图纸 %d 更新成功\n", drawing.ID)
		}
	}

	// 如果没有图纸记录，创建一些新的测试图纸
	if len(drawings) == 0 {
		fmt.Println("没有现有图纸记录，创建新的测试图纸...")

		for i := 1; i <= 3; i++ {
			testDrawing := system.SysDrawing{
				AlbumID:        album.ID,
				SerialNumber:   fmt.Sprintf("TEST-%03d", i),
				Name:           fmt.Sprintf("测试图纸%d", i),
				BeanQuantity:   &[]int{i * 100}[0],
				PosterImageURL: fmt.Sprintf("uploads/test/poster%d.jpg", i),
				DrawingURLs:    fmt.Sprintf(`["uploads/test/drawing%d.pdf", "uploads/test/drawing%d.dwg"]`, i, i),
				CreatorUUID:    user.UUID,
				AllowedMembers: fmt.Sprintf(`["%s"]`, user.UUID.String()),
			}

			err = global.GVA_DB.Create(&testDrawing).Error
			if err != nil {
				log.Printf("创建测试图纸 %d 失败: %v", i, err)
			} else {
				fmt.Printf("成功创建测试图纸 %d，ID: %d\n", i, testDrawing.ID)
			}
		}
	}

	fmt.Println("测试图纸数据创建/更新完成！")
}
