package system

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"gorm.io/gorm"
)

type MustReadService struct{}

// CreateMustRead 创建必读内容
func (mustReadService *MustReadService) CreateMustRead(mustReadReq request.CreateMustRead) (mustRead system.SysMustRead, err error) {
	mustRead = system.SysMustRead{
		CreatorUUID: mustReadReq.CreatorUUID,
		Title:       mustReadReq.Title,
		Content:     mustReadReq.Content,
		Status:      1, // 默认状态为正常
	}
	err = global.GVA_DB.Create(&mustRead).Error
	return mustRead, err
}

// DeleteMustRead 删除必读内容
func (mustReadService *MustReadService) DeleteMustRead(mustReadReq request.DeleteMustRead) error {
	return global.GVA_DB.Delete(&system.SysMustRead{}, mustReadReq.ID).Error
}

// UpdateMustRead 更新必读内容
func (mustReadService *MustReadService) UpdateMustRead(mustReadReq request.UpdateMustRead) error {
	var mustRead system.SysMustRead
	err := global.GVA_DB.Where("id = ?", mustReadReq.ID).First(&mustRead).Error
	if err != nil {
		return err
	}

	// 更新字段
	if mustReadReq.Title != "" {
		mustRead.Title = mustReadReq.Title
	}
	if mustReadReq.Content != "" {
		mustRead.Content = mustReadReq.Content
	}
	if mustReadReq.Status != 0 {
		mustRead.Status = mustReadReq.Status
	}

	return global.GVA_DB.Save(&mustRead).Error
}

// GetMustRead 根据ID获取必读内容
func (mustReadService *MustReadService) GetMustRead(mustReadReq request.GetMustReadByID) (mustRead system.SysMustRead, err error) {
	err = global.GVA_DB.Preload("Creator").Where("id = ?", mustReadReq.ID).First(&mustRead).Error
	return mustRead, err
}

// GetLatestMustRead 获取最新的必读内容
func (mustReadService *MustReadService) GetLatestMustRead() (mustRead system.SysMustRead, err error) {
	err = global.GVA_DB.Preload("Creator").Where("status = ?", 1).Order("created_at desc").First(&mustRead).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，返回空结构体
			return system.SysMustRead{}, nil
		}
		return mustRead, err
	}
	return mustRead, nil
}