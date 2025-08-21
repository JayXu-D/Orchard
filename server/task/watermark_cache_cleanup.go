package task

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/watermark"
	"go.uber.org/zap"
)

// WatermarkCacheCleanupTask 水印缓存清理任务
type WatermarkCacheCleanupTask struct {
	watermarkService *watermark.WatermarkService
}

// NewWatermarkCacheCleanupTask 创建水印缓存清理任务
func NewWatermarkCacheCleanupTask() *WatermarkCacheCleanupTask {
	return &WatermarkCacheCleanupTask{
		watermarkService: watermark.NewWatermarkService(),
	}
}

// Run 运行清理任务
func (t *WatermarkCacheCleanupTask) Run() {
	global.GVA_LOG.Info("开始清理水印缓存...")

	// 清理过期缓存
	if err := t.watermarkService.CleanExpiredCache(); err != nil {
		global.GVA_LOG.Error("清理水印缓存失败", zap.Error(err))
		return
	}

	// 获取缓存大小
	cacheSize, err := t.watermarkService.GetCacheSize()
	if err != nil {
		global.GVA_LOG.Error("获取水印缓存大小失败", zap.Error(err))
		return
	}

	global.GVA_LOG.Info("水印缓存清理完成",
		zap.Int64("cacheSize", cacheSize),
		zap.String("unit", "bytes"))
}

// GetInterval 获取执行间隔
func (t *WatermarkCacheCleanupTask) GetInterval() time.Duration {
	return 6 * time.Hour // 每6小时执行一次
}

// GetName 获取任务名称
func (t *WatermarkCacheCleanupTask) GetName() string {
	return "WatermarkCacheCleanup"
}
