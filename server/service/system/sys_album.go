package system

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	albumRequest "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/google/uuid"
)

type AlbumService struct{}

// CreateAlbum 创建相册
func (albumService *AlbumService) CreateAlbum(albumReq albumRequest.CreateAlbum) (album system.SysAlbum, err error) {
	// 验证创建者是否存在
	var creator system.SysUser
	if err := global.GVA_DB.Where("uuid = ?", albumReq.CreatorUUID).First(&creator).Error; err != nil {
		return album, errors.New("创建者不存在")
	}

	// 创建相册
	album = system.SysAlbum{
		CreatorUUID:   albumReq.CreatorUUID,
		Title:         albumReq.Title,
		CoverImageURL: albumReq.CoverImageURL,
		Description:   albumReq.Description,
		Status:        1, // 默认状态为正常
	}

	// 开启事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存相册
	if err := tx.Create(&album).Error; err != nil {
		tx.Rollback()
		return album, err
	}

	// 如果有管理员，创建关联关系
	if len(albumReq.AdminUserIDs) > 0 {
		// 验证管理员用户是否存在
		var adminUsers []system.SysUser
		if err := tx.Where("id IN ?", albumReq.AdminUserIDs).Find(&adminUsers).Error; err != nil {
			tx.Rollback()
			return album, err
		}

		if len(adminUsers) != len(albumReq.AdminUserIDs) {
			tx.Rollback()
			return album, errors.New("部分管理员用户不存在")
		}

		// 创建相册管理员关联
		for _, adminID := range albumReq.AdminUserIDs {
			albumAdmin := system.SysAlbumAdmin{
				AlbumID: album.ID,
				UserID:  adminID,
			}
			if err := tx.Create(&albumAdmin).Error; err != nil {
				tx.Rollback()
				return album, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return album, err
	}

	// 重新查询相册信息（包含关联数据）
	err = global.GVA_DB.Preload("Creator").Preload("AdminUsers").First(&album, album.ID).Error
	return album, err
}

// DeleteAlbum 删除相册
func (albumService *AlbumService) DeleteAlbum(albumReq albumRequest.DeleteAlbum) (err error) {
	// 开启事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除相册管理员关联
	if err := tx.Where("album_id = ?", albumReq.ID).Delete(&system.SysAlbumAdmin{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除相册
	if err := tx.Delete(&system.SysAlbum{}, albumReq.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// UpdateAlbum 更新相册
func (albumService *AlbumService) UpdateAlbum(albumReq albumRequest.UpdateAlbum) (err error) {
	// 开启事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新相册基本信息
	updateData := map[string]interface{}{
		"title":           albumReq.Title,
		"cover_image_url": albumReq.CoverImageURL,
		"description":     albumReq.Description,
		"status":          albumReq.Status,
	}

	if err := tx.Model(&system.SysAlbum{}).Where("id = ?", albumReq.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果提供了管理员列表，更新管理员关联
	if albumReq.AdminUserIDs != nil {
		// 删除原有关联
		if err := tx.Where("album_id = ?", albumReq.ID).Delete(&system.SysAlbumAdmin{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 创建新关联
		if len(albumReq.AdminUserIDs) > 0 {
			// 验证管理员用户是否存在
			var adminUsers []system.SysUser
			if err := tx.Where("id IN ?", albumReq.AdminUserIDs).Find(&adminUsers).Error; err != nil {
				tx.Rollback()
				return err
			}

			if len(adminUsers) != len(albumReq.AdminUserIDs) {
				tx.Rollback()
				return errors.New("部分管理员用户不存在")
			}

			// 创建新的相册管理员关联
			for _, adminID := range albumReq.AdminUserIDs {
				albumAdmin := system.SysAlbumAdmin{
					AlbumID: albumReq.ID,
					UserID:  adminID,
				}
				if err := tx.Create(&albumAdmin).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// GetAlbum 根据ID获取相册
func (albumService *AlbumService) GetAlbum(albumReq albumRequest.GetAlbumByID) (album system.SysAlbum, err error) {
	err = global.GVA_DB.Preload("Creator").Preload("AdminUsers").First(&album, albumReq.ID).Error
	return album, err
}

// GetAlbumList 获取相册列表
func (albumService *AlbumService) GetAlbumList(info albumRequest.GetAlbumList) (list []system.SysAlbum, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysAlbum{})

	// 添加查询条件
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	// if info.CreatorUUID != uuid.Nil {
	// 	db = db.Where("creator_uuid = ?", info.CreatorUUID)
	// }
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取列表
	err = db.Limit(limit).Offset(offset).Preload("Creator").Preload("AdminUsers").Find(&list).Error
	return list, total, err
}

// GetAlbumsByCreator 根据创建者UUID获取相册列表
func (albumService *AlbumService) GetAlbumsByCreator(creatorUUID uuid.UUID) (list []system.SysAlbum, err error) {
	err = global.GVA_DB.Where("creator_uuid = ?", creatorUUID).Preload("Creator").Preload("AdminUsers").Find(&list).Error
	return list, err
}

// GetAlbumsByAdmin 根据管理员ID获取相册列表
func (albumService *AlbumService) GetAlbumsByAdmin(adminID uint) (list []system.SysAlbum, err error) {
	err = global.GVA_DB.Joins("JOIN sys_album_admin ON sys_albums.id = sys_album_admin.album_id").
		Where("sys_album_admin.user_id = ?", adminID).
		Preload("Creator").Preload("AdminUsers").
		Find(&list).Error
	return list, err
}
