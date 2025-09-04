package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/watermark"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

// GetMyDrawings 获取当前用户可下载的图纸列表
func (drawingService *DrawingService) GetMyDrawings(req request.GetMyDrawings) ([]*system.SysDrawing, int64, error) {
	var drawings []*system.SysDrawing
	var total int64

	// 添加调试日志
	global.GVA_LOG.Info("GetMyDrawings 开始查询",
		zap.String("user_uuid", req.UserUUID),
		zap.Uint("user_id", req.UserID),
		zap.String("keyword", req.Keyword))

	// 构建查询，获取用户有权限下载的图纸
	// 用户有权限下载的图纸包括：
	// 1. 用户创建的图纸
	// 2. 用户是相册管理员的相册中的图纸
	// 3. 图纸的AllowedMembers中包含该用户
	db := global.GVA_DB.Model(&system.SysDrawing{}).
		Joins("LEFT JOIN sys_albums ON sys_drawings.album_id = sys_albums.id").
		Joins("LEFT JOIN sys_album_admin ON sys_albums.id = sys_album_admin.album_id").
		Where("sys_drawings.creator_uuid = ? OR sys_album_admin.user_id = ? OR JSON_CONTAINS(sys_drawings.allowed_members, ?)",
			req.UserUUID, req.UserID, fmt.Sprintf("\"%s\"", req.UserUUID))

	// 添加搜索条件
	if req.Keyword != "" {
		db = db.Where("sys_drawings.serial_number LIKE ? OR sys_drawings.name LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	err := db.Count(&total).Error
	if err != nil {
		global.GVA_LOG.Error("获取总数失败", zap.Error(err))
		return nil, 0, err
	}

	global.GVA_LOG.Info("查询到的总数", zap.Int64("total", total))

	// 分页查询
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		db = db.Offset(offset).Limit(req.PageSize)
	}

	// 预加载关联数据并去重，使用子查询来避免DISTINCT和ORDER BY的冲突
	err = db.Preload("Album").Preload("Creator").
		Order("sys_drawings.created_at DESC").
		Find(&drawings).Error
	if err != nil {
		global.GVA_LOG.Error("查询图纸列表失败", zap.Error(err))
		return nil, 0, err
	}

	global.GVA_LOG.Info("查询完成",
		zap.Int("drawings_count", len(drawings)),
		zap.Int64("total", total))

	// 添加详细的调试日志
	for i, drawing := range drawings {
		global.GVA_LOG.Info("图纸详情",
			zap.Int("index", i),
			zap.Uint("id", drawing.ID),
			zap.String("name", drawing.Name),
			zap.String("serial_number", drawing.SerialNumber),
			zap.Uint("album_id", drawing.AlbumID),
			zap.String("creator_uuid", drawing.CreatorUUID.String()),
			zap.String("album_title", drawing.Album.Title),
			zap.String("creator_username", drawing.Creator.Username))
	}

	return drawings, total, nil
}

// UpdateEmptyDrawings 更新空白图纸记录（临时方法，用于测试）
func (drawingService *DrawingService) UpdateEmptyDrawings() error {
	// 查找所有空白的图纸记录
	var emptyDrawings []system.SysDrawing
	err := global.GVA_DB.Where("name = '' OR name IS NULL OR album_id = 0").Find(&emptyDrawings).Error
	if err != nil {
		return err
	}

	if len(emptyDrawings) == 0 {
		global.GVA_LOG.Info("没有找到空白图纸记录")
		return nil
	}

	global.GVA_LOG.Info("找到空白图纸记录", zap.Int("count", len(emptyDrawings)))

	// 查找第一个可用的相册和用户
	var album system.SysAlbum
	err = global.GVA_DB.First(&album).Error
	if err != nil {
		global.GVA_LOG.Warn("没有找到相册数据")
		return err
	}

	var user system.SysUser
	err = global.GVA_DB.First(&user).Error
	if err != nil {
		global.GVA_LOG.Warn("没有找到用户数据")
		return err
	}

	global.GVA_LOG.Info("使用相册和用户",
		zap.Uint("album_id", album.ID),
		zap.String("album_title", album.Title),
		zap.String("user_uuid", user.UUID.String()),
		zap.String("username", user.Username))

	// 更新空白图纸记录
	for i, drawing := range emptyDrawings {
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
			global.GVA_LOG.Error("更新图纸失败",
				zap.Uint("drawing_id", drawing.ID),
				zap.Error(err))
		} else {
			global.GVA_LOG.Info("图纸更新成功",
				zap.Uint("drawing_id", drawing.ID),
				zap.String("name", updates["name"].(string)))
		}
	}

	return nil
}

// DownloadDrawing 下载图纸
func (drawingService *DrawingService) DownloadDrawing(req request.DownloadDrawing, userUUID uuid.UUID) (*systemRes.DownloadResponse, error) {
	// 获取图纸信息
	var drawing system.SysDrawing
	err := global.GVA_DB.First(&drawing, req.DrawingID).Error
	if err != nil {
		return nil, err
	}

	// 添加调试日志
	global.GVA_LOG.Info("下载图纸",
		zap.Uint("drawing_id", req.DrawingID),
		zap.String("drawing_name", drawing.Name),
		zap.String("drawing_urls", drawing.DrawingURLs))

	// 检查权限（这里可以添加更复杂的权限检查逻辑）
	// TODO: 实现权限检查

	// 记录下载历史
	downloadHistoryService := &DownloadHistoryService{}
	err = downloadHistoryService.RecordDownload(userUUID, req.DrawingID, req.AlbumID)
	if err != nil {
		global.GVA_LOG.Warn("记录下载历史失败", zap.Error(err))
		// 不因为记录失败而阻止下载
	}

	// 解析图纸文件URLs
	var drawingURLs []string
	if drawing.DrawingURLs != "" {
		err = json.Unmarshal([]byte(drawing.DrawingURLs), &drawingURLs)
		if err != nil {
			global.GVA_LOG.Warn("解析图纸URLs失败",
				zap.Uint("drawing_id", req.DrawingID),
				zap.String("drawing_urls", drawing.DrawingURLs),
				zap.Error(err))
			return nil, err
		}

		global.GVA_LOG.Info("解析后的图纸URLs",
			zap.Uint("drawing_id", req.DrawingID),
			zap.Int("urls_count", len(drawingURLs)),
			zap.Any("urls", drawingURLs))
	} else {
		global.GVA_LOG.Warn("图纸没有DrawingURLs",
			zap.Uint("drawing_id", req.DrawingID),
			zap.String("drawing_name", drawing.Name))
	}

	// 获取创建者信息用于水印
	var creator system.SysUser
	if err := global.GVA_DB.Where("uuid = ?", drawing.CreatorUUID).First(&creator).Error; err != nil {
		global.GVA_LOG.Warn("获取创建者信息失败", zap.Error(err))
	}

	// 处理水印
	var filePaths []string
	if req.AddWatermark && len(drawingURLs) > 0 {
		watermarkService := watermark.NewWatermarkService()
		watermarkText := req.WatermarkText
		if watermarkText == "" {
			watermarkText = fmt.Sprintf("创建者: %s", creator.Username)
		}

		// 为每个图纸文件添加水印
		for _, drawingURL := range drawingURLs {
			// 构建完整的文件路径
			// 检查drawingURL是否已经包含uploads前缀
			var fullPath string
			if strings.HasPrefix(drawingURL, "uploads/") {
				fullPath = drawingURL
			} else {
				fullPath = filepath.Join("uploads", drawingURL)
			}

			global.GVA_LOG.Info("检查文件路径",
				zap.String("drawing_url", drawingURL),
				zap.String("full_path", fullPath))

			if _, err := os.Stat(fullPath); err == nil {
				// 添加水印
				watermarkedPath, err := watermarkService.AddWatermark(fullPath, watermarkText)
				if err == nil {
					// 返回可以通过HTTP访问的路径
					httpPath := "/api/v1/drawing/watermark/" + filepath.Base(watermarkedPath)
					filePaths = append(filePaths, httpPath)
					global.GVA_LOG.Info("添加水印成功", zap.String("file", fullPath), zap.String("http_path", httpPath))
				} else {
					global.GVA_LOG.Warn("添加水印失败", zap.String("file", fullPath), zap.Error(err))
					// 水印失败时，返回原文件的HTTP路径
					httpPath := "/api/v1/drawing/file/" + filepath.Base(fullPath)
					filePaths = append(filePaths, httpPath)
				}
			} else {
				global.GVA_LOG.Warn("文件不存在", zap.String("file", fullPath), zap.Error(err))
			}
		}
	} else {
		// 不添加水印，直接使用原文件
		for _, drawingURL := range drawingURLs {
			// 检查drawingURL是否已经包含uploads前缀
			var fullPath string
			if strings.HasPrefix(drawingURL, "uploads/") {
				fullPath = drawingURL
			} else {
				fullPath = filepath.Join("uploads", drawingURL)
			}

			global.GVA_LOG.Info("检查文件路径（无水印）",
				zap.String("drawing_url", drawingURL),
				zap.String("full_path", fullPath))

			if _, err := os.Stat(fullPath); err == nil {
				filePaths = append(filePaths, fullPath)
				global.GVA_LOG.Info("文件存在（无水印）", zap.String("file", fullPath))
			} else {
				global.GVA_LOG.Warn("文件不存在（无水印）", zap.String("file", fullPath), zap.Error(err))
			}
		}
	}

	// 添加最终调试日志
	global.GVA_LOG.Info("下载完成",
		zap.Uint("drawing_id", req.DrawingID),
		zap.Int("total_file_paths", len(filePaths)),
		zap.Any("file_paths", filePaths))

	// 生成下载链接
	fileName := drawing.Name + "_图纸.zip"
	if req.AddWatermark {
		fileName = drawing.Name + "_图纸_水印.zip"
	}

	return &systemRes.DownloadResponse{
		DownloadURL: "/api/v1/drawing/download/" + fmt.Sprintf("%d", req.DrawingID),
		FileName:    fileName,
		FileSize:    0, // TODO: 计算实际文件大小
		FilePaths:   filePaths,
	}, nil
}

// RecordDownload 点击下载时记录下载历史（不返回文件）
func (drawingService *DrawingService) RecordDownload(req request.RecordDownload, userUUID uuid.UUID) error {
	downloadHistoryService := &DownloadHistoryService{}
	return downloadHistoryService.RecordDownload(userUUID, req.DrawingID, req.AlbumID)
}

// BatchDownloadDrawings 批量下载图纸
func (drawingService *DrawingService) BatchDownloadDrawings(req request.BatchDownloadDrawings, userUUID uuid.UUID) (*systemRes.DownloadResponse, error) {
	// 获取所有图纸信息
	var drawings []system.SysDrawing
	err := global.GVA_DB.Where("id IN ?", req.DrawingIDs).Find(&drawings).Error
	if err != nil {
		return nil, err
	}

	// 添加调试日志
	global.GVA_LOG.Info("批量下载图纸",
		zap.Int("requested_count", len(req.DrawingIDs)),
		zap.Int("found_drawings", len(drawings)),
		zap.Any("drawing_ids", req.DrawingIDs))

	// 检查权限（这里可以添加更复杂的权限检查逻辑）
	// TODO: 实现权限检查

	// 记录批量下载历史
	downloadHistoryService := &DownloadHistoryService{}
	for _, drawing := range drawings {
		err = downloadHistoryService.RecordDownload(userUUID, drawing.ID, req.AlbumID)
		if err != nil {
			global.GVA_LOG.Warn("记录下载历史失败", zap.Error(err))
			// 不因为记录失败而阻止下载
		}
	}

	// 收集所有图纸文件URLs
	var allFilePaths []string
	if req.AddWatermark {
		watermarkService := watermark.NewWatermarkService()
		watermarkText := req.WatermarkText
		if watermarkText == "" {
			watermarkText = "批量下载图纸"
		}

		for _, drawing := range drawings {
			// 添加调试日志
			global.GVA_LOG.Info("处理图纸",
				zap.Uint("drawing_id", drawing.ID),
				zap.String("drawing_name", drawing.Name),
				zap.String("drawing_urls", drawing.DrawingURLs))

			var drawingURLs []string
			if drawing.DrawingURLs != "" {
				err = json.Unmarshal([]byte(drawing.DrawingURLs), &drawingURLs)
				if err != nil {
					global.GVA_LOG.Warn("解析图纸URLs失败",
						zap.Uint("drawing_id", drawing.ID),
						zap.String("drawing_urls", drawing.DrawingURLs),
						zap.Error(err))
					continue
				}

				global.GVA_LOG.Info("解析后的图纸URLs",
					zap.Uint("drawing_id", drawing.ID),
					zap.Int("urls_count", len(drawingURLs)),
					zap.Any("urls", drawingURLs))

				// 为每个图纸文件添加水印
				for _, drawingURL := range drawingURLs {
					// 检查drawingURL是否已经包含uploads前缀
					var fullPath string
					if strings.HasPrefix(drawingURL, "uploads/") {
						fullPath = drawingURL
					} else {
						fullPath = filepath.Join("uploads", drawingURL)
					}

					global.GVA_LOG.Info("检查文件路径",
						zap.String("drawing_url", drawingURL),
						zap.String("full_path", fullPath))

					if _, err := os.Stat(fullPath); err == nil {
						watermarkedPath, err := watermarkService.AddWatermark(fullPath, watermarkText)
						if err == nil {
							// 返回可以通过HTTP访问的路径
							httpPath := "/api/v1/drawing/watermark/" + filepath.Base(watermarkedPath)
							allFilePaths = append(allFilePaths, httpPath)
							global.GVA_LOG.Info("添加水印成功", zap.String("file", fullPath), zap.String("http_path", httpPath))
						} else {
							global.GVA_LOG.Warn("添加水印失败", zap.String("file", fullPath), zap.Error(err))
							// 水印失败时，返回原文件的HTTP路径
							httpPath := "/api/v1/drawing/file/" + filepath.Base(fullPath)
							allFilePaths = append(allFilePaths, httpPath)
						}
					} else {
						global.GVA_LOG.Warn("文件不存在", zap.String("file", fullPath), zap.Error(err))
					}
				}
			} else {
				global.GVA_LOG.Warn("图纸没有DrawingURLs",
					zap.Uint("drawing_id", drawing.ID),
					zap.String("drawing_name", drawing.Name))
			}
		}
	} else {
		// 不添加水印，直接收集原文件路径
		for _, drawing := range drawings {
			// 添加调试日志
			global.GVA_LOG.Info("处理图纸（无水印）",
				zap.Uint("drawing_id", drawing.ID),
				zap.String("drawing_name", drawing.Name),
				zap.String("drawing_urls", drawing.DrawingURLs))

			var drawingURLs []string
			if drawing.DrawingURLs != "" {
				err = json.Unmarshal([]byte(drawing.DrawingURLs), &drawingURLs)
				if err != nil {
					global.GVA_LOG.Warn("解析图纸URLs失败",
						zap.Uint("drawing_id", drawing.ID),
						zap.String("drawing_urls", drawing.DrawingURLs),
						zap.Error(err))
					continue
				}

				global.GVA_LOG.Info("解析后的图纸URLs（无水印）",
					zap.Uint("drawing_id", drawing.ID),
					zap.Int("urls_count", len(drawingURLs)),
					zap.Any("urls", drawingURLs))

				for _, drawingURL := range drawingURLs {
					// 检查drawingURL是否已经包含uploads前缀
					var fullPath string
					if strings.HasPrefix(drawingURL, "uploads/") {
						fullPath = drawingURL
					} else {
						fullPath = filepath.Join("uploads", drawingURL)
					}

					global.GVA_LOG.Info("检查文件路径（无水印）",
						zap.String("drawing_url", drawingURL),
						zap.String("full_path", fullPath))

					if _, err := os.Stat(fullPath); err == nil {
						allFilePaths = append(allFilePaths, fullPath)
						global.GVA_LOG.Info("文件存在（无水印）", zap.String("file", fullPath))
					} else {
						global.GVA_LOG.Warn("文件不存在（无水印）", zap.String("file", fullPath), zap.Error(err))
					}
				}
			} else {
				global.GVA_LOG.Warn("图纸没有DrawingURLs（无水印）",
					zap.Uint("drawing_id", drawing.ID),
					zap.String("drawing_name", drawing.Name))
			}
		}
	}

	// 添加最终调试日志
	global.GVA_LOG.Info("批量下载完成",
		zap.Int("total_file_paths", len(allFilePaths)),
		zap.Any("file_paths", allFilePaths))

	// 生成批量下载链接
	fileName := "批量图纸_" + time.Now().Format("2006-01-02") + ".zip"
	if req.AddWatermark {
		fileName = "批量图纸_水印_" + time.Now().Format("2006-01-02") + ".zip"
	}

	return &systemRes.DownloadResponse{
		DownloadURL: "/api/v1/drawing/batchDownload/" + fmt.Sprintf("%d", req.AlbumID),
		FileName:    fileName,
		FileSize:    0, // TODO: 计算实际文件大小
		FilePaths:   allFilePaths,
	}, nil
}
