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
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
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

// AddWatermark 为图片添加文字水印（参考博文方法：小图文字->旋转->平铺）
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

	base := image.NewRGBA(img.Bounds())
	draw.Draw(base, base.Bounds(), img, image.Point{}, draw.Src)

	// 依据图片尺寸设定字号
	imgW := base.Bounds().Dx()
	imgH := base.Bounds().Dy()
	fontSize := math.Max(18, float64(minInt(imgW, imgH))/40)
	// 生成仅包含文字的透明图
	wmColor := color.RGBA{R: 255, G: 255, B: 255, A: 100}
	wmImg, err := MakeImageByText(watermarkText, wmColor, color.Transparent, fontSize)
	if err != nil {
		return "", err
	}
	// 旋转小图（-30°），透明背景
	rotated := imaging.Rotate(wmImg, 30, color.Transparent)

	// 平铺覆盖整张图，基于图片尺寸自适应间距
	cols := 5
	rows := 5
	tileSpacingX := maxInt(rotated.Bounds().Dx(), imgW/cols)
	tileSpacingY := maxInt(rotated.Bounds().Dy(), imgH/rows)

	startX := -rotated.Bounds().Dx()
	startY := -rotated.Bounds().Dy()
	for y := startY; y < imgH+rotated.Bounds().Dy(); y += tileSpacingY {
		rowOffset := 0
		if ((y-startY)/tileSpacingY)%2 == 1 {
			rowOffset = tileSpacingX / 2
		}
		for x := startX + rowOffset; x < imgW+rotated.Bounds().Dx(); x += tileSpacingX {
			r := rotated.Bounds().Add(image.Pt(x, y))
			draw.Draw(base, r, rotated, rotated.Bounds().Min, draw.Over)
		}
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
		err = png.Encode(out, base)
	case "jpeg", "jpg":
		// 使用 4:4:4 编码尽量降低伪影
		ycc := toYCbCr444(base)
		err = jpeg.Encode(out, ycc, &jpeg.Options{Quality: 100})
	default:
		ycc := toYCbCr444(base)
		err = jpeg.Encode(out, ycc, &jpeg.Options{Quality: 100})
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
	// 必须是非空文件
	if info.Size() == 0 {
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

// minInt returns the smaller of two ints
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// renderLabelImage renders multiline text into a transparent RGBA image using the provided face and color
func renderLabelImage(text string, face font.Face, col color.RGBA) *image.RGBA {
	lines := strings.Split(text, "\n")
	for i := range lines {
		if len(lines[i]) > 64 {
			lines[i] = lines[i][:64]
		}
	}

	dr := &font.Drawer{Face: face}
	maxAdvance := dr.MeasureString(longestLine(lines)).Ceil()

	metrics := face.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	if lineHeight <= 0 {
		lineHeight = 16
	}

	padding := 12
	width := maxAdvance + padding*2
	height := lineHeight*len(lines) + padding*2

	// 1. 初始化目标图像为透明背景（RGBA: 0,0,0,0）
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// 2. 填充背景为透明（可选，若需纯色背景可替换为其他颜色）
	// 例如：填充白色背景
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 0}}, image.Point{}, draw.Src)

	dr.Dst = rgba
	// 3. 设置文字颜色（包含 Alpha 通道）
	dr.Src = image.NewUniform(col)
	dr.Face = face

	x := padding
	y := padding + lineHeight
	for _, line := range lines {
		dr.Dot = fixed.P(x, y)
		dr.DrawString(line)
		y += lineHeight
	}
	return rgba
}

// renderLabelMask 将文本绘制为 Alpha 掩码，避免半透明边缘噪点
func renderLabelMask(text string, face font.Face) *image.Alpha {
	lines := strings.Split(text, "\n")
	for i := range lines {
		if len(lines[i]) > 64 {
			lines[i] = lines[i][:64]
		}
	}

	dr := &font.Drawer{Face: face}
	maxAdvance := dr.MeasureString(longestLine(lines)).Ceil()
	metrics := face.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	if lineHeight <= 0 {
		lineHeight = 16
	}
	padding := 12
	width := maxAdvance + padding*2
	height := lineHeight*len(lines) + padding*2

	alpha := image.NewAlpha(image.Rect(0, 0, width, height))
	dr.Dst = alpha
	dr.Src = image.White
	dr.Face = face

	x := padding
	y := padding + lineHeight
	for _, line := range lines {
		dr.Dot = fixed.P(x, y)
		dr.DrawString(line)
		y += lineHeight
	}
	return alpha
}

// rotateAlpha 以双线性插值旋转 Alpha 掩码，减少锯齿与噪点
func rotateAlpha(src *image.Alpha, angleRad float64) *image.Alpha {
	sw := src.Bounds().Dx()
	sh := src.Bounds().Dy()
	cx := float64(sw) / 2.0
	cy := float64(sh) / 2.0

	corners := [][2]float64{{0, 0}, {float64(sw), 0}, {0, float64(sh)}, {float64(sw), float64(sh)}}
	sinA := math.Sin(angleRad)
	cosA := math.Cos(angleRad)
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, c := range corners {
		x := c[0] - cx
		y := c[1] - cy
		rx := x*cosA - y*sinA
		ry := x*sinA + y*cosA
		rx += cx
		ry += cy
		if rx < minX {
			minX = rx
		}
		if ry < minY {
			minY = ry
		}
		if rx > maxX {
			maxX = rx
		}
		if ry > maxY {
			maxY = ry
		}
	}

	dw := int(math.Ceil(maxX - minX))
	dh := int(math.Ceil(maxY - minY))
	dst := image.NewAlpha(image.Rect(0, 0, dw, dh))

	dcx := float64(dw) / 2.0
	dcy := float64(dh) / 2.0

	for dy := 0; dy < dh; dy++ {
		for dx := 0; dx < dw; dx++ {
			x := float64(dx) - dcx
			y := float64(dy) - dcy
			sx := x*cosA + y*sinA + cx
			sy := -x*sinA + y*cosA + cy

			if sx < 0 || sy < 0 || sx > float64(sw-1) || sy > float64(sh-1) {
				continue
			}

			x0 := int(math.Floor(sx))
			y0 := int(math.Floor(sy))
			x1 := x0 + 1
			y1 := y0 + 1
			fx := sx - float64(x0)
			fy := sy - float64(y0)
			if x1 >= sw {
				x1 = sw - 1
			}
			if y1 >= sh {
				y1 = sh - 1
			}

			o00 := src.PixOffset(x0, y0)
			o10 := src.PixOffset(x1, y0)
			o01 := src.PixOffset(x0, y1)
			o11 := src.PixOffset(x1, y1)

			a00 := float64(src.Pix[o00])
			a10 := float64(src.Pix[o10])
			a01 := float64(src.Pix[o01])
			a11 := float64(src.Pix[o11])
			wa := (1 - fx) * (1 - fy)
			wb := fx * (1 - fy)
			wc := (1 - fx) * fy
			wd := fx * fy
			a := a00*wa + a10*wb + a01*wc + a11*wd
			if a <= 0.5 {
				continue
			}
			dOff := dst.PixOffset(dx, dy)
			dst.Pix[dOff] = uint8(math.Min(255, math.Max(0, a)))
		}
	}

	return dst
}

// rotateRGBA 以双线性插值旋转 RGBA 图像，减少锯齿与噪点
func rotateRGBA(src *image.RGBA, angleRad float64) *image.RGBA {
	sw := src.Bounds().Dx()
	sh := src.Bounds().Dy()
	cx := float64(sw) / 2.0
	cy := float64(sh) / 2.0

	corners := [][2]float64{{0, 0}, {float64(sw), 0}, {0, float64(sh)}, {float64(sw), float64(sh)}}
	sinA := math.Sin(angleRad)
	cosA := math.Cos(angleRad)
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, c := range corners {
		x := c[0] - cx
		y := c[1] - cy
		rx := x*cosA - y*sinA
		ry := x*sinA + y*cosA
		rx += cx
		ry += cy
		if rx < minX {
			minX = rx
		}
		if ry < minY {
			minY = ry
		}
		if rx > maxX {
			maxX = rx
		}
		if ry > maxY {
			maxY = ry
		}
	}

	dw := int(math.Ceil(maxX - minX))
	dh := int(math.Ceil(maxY - minY))
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))

	dcx := float64(dw) / 2.0
	dcy := float64(dh) / 2.0

	for dy := 0; dy < dh; dy++ {
		for dx := 0; dx < dw; dx++ {
			x := float64(dx) - dcx
			y := float64(dy) - dcy
			sx := x*cosA + y*sinA + cx
			sy := -x*sinA + y*cosA + cy

			if sx < 0 || sy < 0 || sx > float64(sw-1) || sy > float64(sh-1) {
				continue
			}

			x0 := int(math.Floor(sx))
			y0 := int(math.Floor(sy))
			x1 := x0 + 1
			y1 := y0 + 1
			fx := sx - float64(x0)
			fy := sy - float64(y0)
			if x1 >= sw {
				x1 = sw - 1
			}
			if y1 >= sh {
				y1 = sh - 1
			}

			o00 := src.PixOffset(x0, y0)
			o10 := src.PixOffset(x1, y0)
			o01 := src.PixOffset(x0, y1)
			o11 := src.PixOffset(x1, y1)

			r00 := float64(src.Pix[o00+0])
			r10 := float64(src.Pix[o10+0])
			r01 := float64(src.Pix[o01+0])
			r11 := float64(src.Pix[o11+0])
			g00 := float64(src.Pix[o00+1])
			g10 := float64(src.Pix[o10+1])
			g01 := float64(src.Pix[o01+1])
			g11 := float64(src.Pix[o11+1])
			b00 := float64(src.Pix[o00+2])
			b10 := float64(src.Pix[o10+2])
			b01 := float64(src.Pix[o01+2])
			b11 := float64(src.Pix[o11+2])
			wa := (1 - fx) * (1 - fy)
			wb := fx * (1 - fy)
			wc := (1 - fx) * fy
			wd := fx * fy
			r := r00*wa + r10*wb + r01*wc + r11*wd
			g := g00*wa + g10*wb + g01*wc + g11*wd
			b := b00*wa + b10*wb + b01*wc + b11*wd
			dst.Pix[dst.PixOffset(dx, dy)+0] = uint8(math.Min(255, math.Max(0, r)))
			dst.Pix[dst.PixOffset(dx, dy)+1] = uint8(math.Min(255, math.Max(0, g)))
			dst.Pix[dst.PixOffset(dx, dy)+2] = uint8(math.Min(255, math.Max(0, b)))
		}
	}

	return dst
}

// toYCbCr444 converts an RGBA image to a 4:4:4 YCbCr image to avoid chroma subsampling artifacts when encoding JPEG
func toYCbCr444(src *image.RGBA) *image.YCbCr {
	b := src.Bounds()
	ycc := image.NewYCbCr(b, image.YCbCrSubsampleRatio444)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			o := src.PixOffset(x, y)
			r := src.Pix[o+0]
			g := src.Pix[o+1]
			bl := src.Pix[o+2]
			Y, Cb, Cr := color.RGBToYCbCr(r, g, bl)
			i := ycc.YOffset(x, y)
			ycc.Y[i] = Y
			ycc.Cb[i] = Cb
			ycc.Cr[i] = Cr
		}
	}
	return ycc
}

// MakeImageByText 根据文本内容制作一个仅包含该文本内容的图片
func MakeImageByText(text string, fontColor color.Color, bgColor color.Color, fontSize float64) (image.Image, error) {
	ftCtx := MakeFreetypeCtx(fontSize, fontColor)
	if ftCtx == nil {
		return nil, errors.New("freetype context init failed")
	}
	// 简单估算宽高：字号 * 文本长度；高度取 2*字号
	w := int(fontSize) * int(maxInt(1, len(text)))
	h := int(fontSize) * 2
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	if bgColor != color.Transparent {
		bg := image.NewUniform(bgColor)
		draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	}
	ftCtx.SetClip(rgba.Rect)
	ftCtx.SetDst(rgba)
	pt := freetype.Pt(0, int(ftCtx.PointToFixed(fontSize)>>6))
	if _, err := ftCtx.DrawString(text, pt); err != nil {
		return nil, err
	}
	return rgba, nil
}

// MustParseFont 解析字体（使用 gofont goregular）
func MustParseFont() *truetype.Font {
	ft, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return ft
}

// MakeFreetypeCtx 初始化 freetype 上下文
func MakeFreetypeCtx(fontSize float64, fontColor color.Color) *freetype.Context {
	ctx := freetype.NewContext()
	ctx.SetDPI(100)
	ctx.SetFont(MustParseFont())
	ctx.SetFontSize(fontSize)
	ctx.SetSrc(image.NewUniform(fontColor))
	ctx.SetHinting(font.HintingNone)
	return ctx
}
