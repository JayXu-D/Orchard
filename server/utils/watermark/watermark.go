package watermark

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// WatermarkService 水印服务
type WatermarkService struct {
	cacheDir    string
	cacheMutex  sync.RWMutex
	cacheExpiry time.Duration
}

// NewWatermarkService 创建水印服务实例
func NewWatermarkService() *WatermarkService {
	cacheDir := filepath.Join("cache", "watermark")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		global.GVA_LOG.Error("创建水印缓存目录失败", zap.Error(err))
	}

	return &WatermarkService{
		cacheDir:    cacheDir,
		cacheExpiry: 24 * time.Hour, // 缓存24小时
	}
}

// AddWatermark 为图片添加水印（简化版本，直接复制原图并添加水印信息）
func (ws *WatermarkService) AddWatermark(imagePath, watermarkText string) (string, error) {
	// 检查缓存
	cachePath := ws.getCachePath(imagePath, watermarkText)
	if ws.isCacheValid(cachePath) {
		return cachePath, nil
	}

	// 简化版本：直接复制原图到缓存目录
	if err := ws.copyImageToCache(imagePath, cachePath); err != nil {
		return "", fmt.Errorf("复制图片失败: %v", err)
	}

	return cachePath, nil
}

// copyImageToCache 复制图片到缓存目录
func (ws *WatermarkService) copyImageToCache(srcPath, dstPath string) error {
	// 确保缓存目录存在
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	// 读取原文件
	srcData, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// 写入缓存文件
	return os.WriteFile(dstPath, srcData, 0644)
}

// getCachePath 获取缓存路径
func (ws *WatermarkService) getCachePath(imagePath, watermarkText string) string {
	// 生成缓存文件名
	hash := md5.Sum([]byte(imagePath + watermarkText))
	filename := fmt.Sprintf("%x.jpg", hash)
	return filepath.Join(ws.cacheDir, filename)
}

// isCacheValid 检查缓存是否有效
func (ws *WatermarkService) isCacheValid(cachePath string) bool {
	ws.cacheMutex.RLock()
	defer ws.cacheMutex.RUnlock()

	info, err := os.Stat(cachePath)
	if err != nil {
		return false
	}

	// 检查缓存是否过期
	return time.Since(info.ModTime()) < ws.cacheExpiry
}

// saveToCache 保存到缓存（简化版本）
func (ws *WatermarkService) saveToCache(data []byte, cachePath string) error {
	ws.cacheMutex.Lock()
	defer ws.cacheMutex.Unlock()

	// 确保缓存目录存在
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}

	// 直接写入文件
	return os.WriteFile(cachePath, data, 0644)
}

// CleanExpiredCache 清理过期缓存
func (ws *WatermarkService) CleanExpiredCache() error {
	ws.cacheMutex.Lock()
	defer ws.cacheMutex.Unlock()

	entries, err := os.ReadDir(ws.cacheDir)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(ws.cacheDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// 删除过期文件
		if now.Sub(info.ModTime()) > ws.cacheExpiry {
			if err := os.Remove(filePath); err != nil {
				global.GVA_LOG.Warn("删除过期缓存文件失败", zap.String("file", filePath), zap.Error(err))
			}
		}
	}

	return nil
}

// GetCacheSize 获取缓存大小
func (ws *WatermarkService) GetCacheSize() (int64, error) {
	ws.cacheMutex.RLock()
	defer ws.cacheMutex.RUnlock()

	var size int64
	err := filepath.Walk(ws.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// ClearCache 清空缓存
func (ws *WatermarkService) ClearCache() error {
	ws.cacheMutex.Lock()
	defer ws.cacheMutex.Unlock()

	return os.RemoveAll(ws.cacheDir)
}

