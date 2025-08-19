package system

import (
	"encoding/json"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"gorm.io/gorm"
)

type DrawingService struct{}

// CreateDrawing 创建图纸
func (drawingService *DrawingService) CreateDrawing(req request.CreateDrawing) (*system.SysDrawing, error) {
	// 检查序号是否已存在
	var existingDrawing system.SysDrawing
	err := global.GVA_DB.Where("album_id = ? AND serial_number = ?", req.AlbumID, req.SerialNumber).First(&existingDrawing).Error
	if err == nil {
		return nil, errors.New("该序号已存在")
	}

	// 将图纸文件URLs转换为JSON字符串
	drawingURLsJSON, err := json.Marshal(req.DrawingURLs)
	if err != nil {
		return nil, err
	}

	// 将允许下载的成员UUIDs转换为JSON字符串
	allowedMembersJSON, err := json.Marshal(req.AllowedMemberUUIDs)
	if err != nil {
		return nil, err
	}

	drawing := &system.SysDrawing{
		AlbumID:        req.AlbumID,
		SerialNumber:   req.SerialNumber,
		Name:           req.Name,
		BeanQuantity:   req.BeanQuantity,
		PosterImageURL: req.PosterImageURL,
		DrawingURLs:    string(drawingURLsJSON),
		CreatorUUID:    req.CreatorUUID,
		AllowedMembers: string(allowedMembersJSON),
	}

	err = global.GVA_DB.Create(drawing).Error
	if err != nil {
		return nil, err
	}

	// 预加载关联数据
	err = global.GVA_DB.Preload("Album").Preload("Creator").First(drawing, drawing.ID).Error
	if err != nil {
		return nil, err
	}

	return drawing, nil
}

// UpdateDrawing 更新图纸
func (drawingService *DrawingService) UpdateDrawing(req request.UpdateDrawing) error {
	// 检查图纸是否存在
	var existingDrawing system.SysDrawing
	err := global.GVA_DB.First(&existingDrawing, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("图纸不存在")
		}
		return err
	}

	// 检查序号是否已被其他图纸使用
	if req.SerialNumber != existingDrawing.SerialNumber {
		var duplicateDrawing system.SysDrawing
		err = global.GVA_DB.Where("album_id = ? AND serial_number = ? AND id != ?", req.AlbumID, req.SerialNumber, req.ID).First(&duplicateDrawing).Error
		if err == nil {
			return errors.New("该序号已被其他图纸使用")
		}
	}

	// 将图纸文件URLs转换为JSON字符串
	drawingURLsJSON, err := json.Marshal(req.DrawingURLs)
	if err != nil {
		return err
	}

	// 将允许下载的成员UUIDs转换为JSON字符串
	allowedMembersJSON, err := json.Marshal(req.AllowedMemberUUIDs)
	if err != nil {
		return err
	}

	// 更新图纸
	updates := map[string]interface{}{
		"album_id":         req.AlbumID,
		"serial_number":    req.SerialNumber,
		"name":             req.Name,
		"bean_quantity":    req.BeanQuantity,
		"poster_image_url": req.PosterImageURL,
		"drawing_urls":     string(drawingURLsJSON),
		"allowed_members":  string(allowedMembersJSON),
	}

	return global.GVA_DB.Model(&existingDrawing).Updates(updates).Error
}

// DeleteDrawing 删除图纸
func (drawingService *DrawingService) DeleteDrawing(req request.DeleteDrawing) error {
	return global.GVA_DB.Delete(&system.SysDrawing{}, req.ID).Error
}

// GetDrawingByID 根据ID获取图纸
func (drawingService *DrawingService) GetDrawingByID(req request.GetDrawingByID) (*system.SysDrawing, error) {
	var drawing system.SysDrawing
	err := global.GVA_DB.Preload("Album").Preload("Creator").First(&drawing, req.ID).Error
	if err != nil {
		return nil, err
	}
	return &drawing, nil
}

// GetDrawingList 获取图纸列表
func (drawingService *DrawingService) GetDrawingList(req request.GetDrawingList) ([]*system.SysDrawing, int64, error) {
	var drawings []*system.SysDrawing
	var total int64

	db := global.GVA_DB.Model(&system.SysDrawing{}).Where("album_id = ?", req.AlbumID)

	// 添加搜索条件
	if req.Keyword != "" {
		db = db.Where("serial_number LIKE ? OR name LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 添加创建者筛选
	if req.CreatorID != 0 {
		db = db.Where("creator_id = ?", req.CreatorID)
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		db = db.Offset(offset).Limit(req.PageSize)
	}

	// 预加载关联数据
	err = db.Preload("Album").Preload("Creator").Order("created_at DESC").Find(&drawings).Error
	if err != nil {
		return nil, 0, err
	}

	return drawings, total, nil
}
