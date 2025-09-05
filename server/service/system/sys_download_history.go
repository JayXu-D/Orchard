package system

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DownloadHistoryService struct{}

// RecordDownload 记录下载历史
func (s *DownloadHistoryService) RecordDownload(userUUID uuid.UUID, drawingID, albumID uint) error {
	history := &system.SysDownloadHistory{
		UserUUID:   userUUID,
		DrawingID:  drawingID,
		AlbumID:    albumID,
		DownloadAt: time.Now().Unix(),
	}

	err := global.GVA_DB.Create(history).Error
	if err != nil {
		global.GVA_LOG.Error("记录下载历史失败", zap.Error(err))
		return err
	}

	return nil
}

// GetUserDrawingDownloadHistory 获取用户图纸下载历史
func (s *DownloadHistoryService) GetUserDrawingDownloadHistory(userUUID uuid.UUID, drawingID uint) (*system.SysDownloadHistory, error) {
	var history system.SysDownloadHistory
	err := global.GVA_DB.Where("user_uuid = ? AND drawing_id = ?", userUUID, drawingID).
		Order("download_at DESC").
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// GetUserDrawingsWithDownloadInfo 获取用户图纸列表（包含下载信息）
func (s *DownloadHistoryService) GetUserDrawingsWithDownloadInfo(userUUID uuid.UUID) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	// 先获取用户ID
	var user system.SysUser
	err := global.GVA_DB.Where("uuid = ?", userUUID).First(&user).Error
	if err != nil {
		global.GVA_LOG.Error("获取用户信息失败", zap.Error(err))
		return results, err
	}

	global.GVA_LOG.Info("查询用户图纸列表",
		zap.String("userUUID", userUUID.String()),
		zap.Uint("userID", user.ID))

	// 简化查询：先获取所有图纸，然后过滤权限
	query := `
		SELECT
			d.id,
			d.serial_number,
			d.name,
			d.album_id,
			a.title AS album_title,
			ldt.last_download_time,
			fdt.first_download_time,
			d.created_at AS created_at
		FROM sys_drawings d
		LEFT JOIN sys_albums a ON d.album_id = a.id
		LEFT JOIN (
			SELECT drawing_id, MAX(download_at) AS last_download_time
			FROM sys_download_histories
			WHERE user_uuid = ?
			GROUP BY drawing_id
		) ldt ON ldt.drawing_id = d.id
		LEFT JOIN (
			SELECT drawing_id, MIN(download_at) AS first_download_time
			FROM sys_download_histories
			WHERE user_uuid = ?
			GROUP BY drawing_id
		) fdt ON fdt.drawing_id = d.id
		WHERE d.creator_uuid = ?
		   OR JSON_CONTAINS(d.allowed_members, ?)
		   OR EXISTS (
			   SELECT 1 FROM sys_album_admin aa
			   WHERE aa.album_id = d.album_id
			   AND aa.user_id = ?
		   )
		ORDER BY d.created_at DESC
	`

	rows, err := global.GVA_DB.Raw(
		query,
		userUUID.String(),           // ldt.user_uuid = ?
		userUUID.String(),           // fdt.user_uuid = ?
		userUUID.String(),           // d.creator_uuid = ?
		"\""+userUUID.String()+"\"", // JSON_CONTAINS(..., ?)
		user.ID,                     // aa.user_id = ?
	).Rows()
	if err != nil {
		global.GVA_LOG.Error("执行查询失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var result map[string]interface{}
		var id uint
		var serialNumber, name string
		var albumID uint
		var albumTitle string
		var lastDownloadTime, firstDownloadTime *int64
		var createdAt time.Time

		err := rows.Scan(&id, &serialNumber, &name, &albumID, &albumTitle,
			&lastDownloadTime, &firstDownloadTime, &createdAt)
		if err != nil {
			global.GVA_LOG.Error("扫描行数据失败", zap.Error(err))
			continue
		}

		// 获取时间信息
		var acquisitionTime, lastDownloadTimeStr string

		// 获取时间（第一次下载时间或图纸创建时间）
		if firstDownloadTime != nil {
			acquisitionTime = time.Unix(*firstDownloadTime, 0).Format("2006-01-02T15:04:05Z")
		} else {
			// 如果没有下载记录，使用图纸创建时间（直接使用已查询的 createdAt）
			acquisitionTime = createdAt.Format("2006-01-02T15:04:05Z")
		}

		// 最后一次下载时间
		if lastDownloadTime != nil {
			lastDownloadTimeStr = time.Unix(*lastDownloadTime, 0).Format("2006-01-02T15:04:05Z")
		}

		result = map[string]interface{}{
			"id":               id,
			"serialNumber":     serialNumber,
			"name":             name,
			"albumId":          albumID,
			"albumTitle":       albumTitle,
			"acquisitionTime":  acquisitionTime,
			"lastDownloadTime": lastDownloadTimeStr,
		}

		results = append(results, result)
		count++
	}

	global.GVA_LOG.Info("查询结果",
		zap.String("userUUID", userUUID.String()),
		zap.Int("resultCount", count),
		zap.Int("totalResults", len(results)))

	return results, nil
}

// GetDownloadStatus 获取指定图纸ID列表的最新下载时间（按用户）
func (s *DownloadHistoryService) GetDownloadStatus(userUUID uuid.UUID, drawingIDs []uint) (map[uint]int64, error) {
	if len(drawingIDs) == 0 {
		return map[uint]int64{}, nil
	}
	type row struct {
		DrawingID uint
		LastTime  *int64
	}
	var rows []row
	err := global.GVA_DB.Table("sys_download_histories").
		Select("drawing_id as drawing_id, MAX(download_at) as last_time").
		Where("user_uuid = ? AND drawing_id IN ?", userUUID, drawingIDs).
		Group("drawing_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint]int64, len(rows))
	for _, r := range rows {
		if r.LastTime != nil {
			result[r.DrawingID] = *r.LastTime
		}
	}
	return result, nil
}
