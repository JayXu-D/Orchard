package watermark

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
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

// AddWatermark 为图片添加文字水印（右下角半透明背景+文本）
func (ws *WatermarkService) AddWatermark(imagePath, watermarkText string) (string, error) {
	cachePath := ws.getCachePath(imagePath, watermarkText)
	if ws.isCacheValid(cachePath) {
		return cachePath, nil
	}

	f, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	img, format, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	// 默认使用 basicfont，作为回退
	face := font.Face(basicfont.Face7x13)
	// 尝试加载 13px TrueType 字体
	if ft, e := opentype.Parse(goregular.TTF); e == nil {
		if f13, e2 := opentype.NewFace(ft, &opentype.FaceOptions{Size: 13, DPI: 72, Hinting: font.HintingFull}); e2 == nil {
			face = f13
		}
	}

	dr := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}), // 纯白，不透明
		Face: face,
	}

	paddingX := 16
	paddingY := 16
	lines := strings.Split(watermarkText, "\n")
	if len(lines) == 1 && len(lines[0]) > 64 {
		lines[0] = lines[0][:64]
	}

	// 行高依据字体大小做简单估算（13px 字体约 16px 行高）
	lineHeight := 16
	if face == basicfont.Face7x13 {
		lineHeight = 14
	}
	totalHeight := lineHeight * len(lines)
	bgHeight := totalHeight + 12
	bgWidth := dr.MeasureString(longestLine(lines)).Ceil() + 12
	bgRect := image.Rect(
		rgba.Bounds().Max.X-bgWidth-paddingX,
		rgba.Bounds().Max.Y-bgHeight-paddingY,
		rgba.Bounds().Max.X-paddingX,
		rgba.Bounds().Max.Y-paddingY,
	)
	// 背景透明度由外部设置或默认值，保持现有
	draw.Draw(rgba, bgRect, image.NewUniform(color.RGBA{0, 0, 0, 50}), image.Point{}, draw.Over)

	x := rgba.Bounds().Max.X - bgWidth - paddingX + 6
	y := rgba.Bounds().Max.Y - bgHeight - paddingY + 6 + lineHeight
	for _, line := range lines {
		dr.Dot = fixed.P(x, y)
		dr.DrawString(line)
		y += lineHeight
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return "", err
	}
	out, err := os.Create(cachePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = out.Close() }()

	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(out, rgba)
	case "jpeg", "jpg":
		err = jpeg.Encode(out, rgba, &jpeg.Options{Quality: 100})
	default:
		err = jpeg.Encode(out, rgba, &jpeg.Options{Quality: 100})
	}
	if err != nil {
		return "", err
	}

	if info, err := os.Stat(cachePath); err == nil && info.Size() == 0 {
		return "", errors.New("watermark output is empty")
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

// longestLine returns the line with maximum visual width for measurement
func longestLine(lines []string) string {
	longest := ""
	for _, l := range lines {
		if len(l) > len(longest) {
			longest = l
		}
	}
	return longest
}
