# 水印功能说明

## 功能概述

本模块为图纸下载提供水印功能，确保下载的图纸包含创建者信息，防止图纸被滥用。

## 主要特性

- **自动水印添加**：下载图纸时自动添加创建者UUID水印
- **智能缓存管理**：缓存已添加水印的图片，避免重复处理
- **定时清理**：自动清理过期的缓存文件，节省存储空间
- **批量处理**：支持批量下载时统一添加水印

## 使用方法

### 1. 单个图纸下载

```javascript
const result = await downloadDrawingApi({
  drawingId: drawing.id,
  albumId: Number(albumId.value),
  addWatermark: true, // 启用水印
  watermarkText: `创建者: ${username}` // 自定义水印文字
})
```

### 2. 批量图纸下载

```javascript
const result = await batchDownloadDrawings({
  drawingIds: selectedDrawings.value,
  albumId: Number(albumId.value),
  addWatermark: true, // 启用水印
  watermarkText: '批量下载图纸' // 自定义水印文字
})
```

## 缓存管理

### 缓存目录
- 默认缓存目录：`cache/watermark/`
- 缓存过期时间：24小时

### 自动清理
- 系统每6小时自动清理过期缓存
- 可通过API手动清理缓存

### 缓存清理API
```go
// 清理过期缓存
watermarkService.CleanExpiredCache()

// 获取缓存大小
cacheSize, err := watermarkService.GetCacheSize()

// 清空所有缓存
watermarkService.ClearCache()
```

## 配置选项

### 水印选项
- `addWatermark`: 是否添加水印（默认true）
- `watermarkText`: 水印文字内容
- 如果不指定水印文字，系统会自动使用"创建者: {用户名}"格式

### 缓存选项
- 缓存目录：可通过环境变量配置
- 缓存过期时间：24小时（可配置）
- 清理间隔：6小时（可配置）

## 注意事项

1. **文件格式支持**：目前支持JPEG、PNG等常见图片格式
2. **水印位置**：水印默认添加在图片右下角
3. **性能考虑**：首次添加水印会有处理时间，后续使用缓存
4. **存储空间**：水印图片会占用额外存储空间，系统会自动清理

## 故障排除

### 常见问题

1. **水印添加失败**
   - 检查图片文件是否存在
   - 检查文件权限
   - 查看系统日志

2. **缓存清理失败**
   - 检查缓存目录权限
   - 确认定时任务是否正常运行

3. **下载失败**
   - 检查文件路径是否正确
   - 确认水印图片是否生成成功

### 日志查看

水印相关的日志会记录在系统日志中，包括：
- 水印添加成功/失败
- 缓存清理状态
- 文件处理错误

## 扩展功能

未来可扩展的功能：
- 支持更多图片格式
- 自定义水印位置和样式
- 水印透明度调节
- 批量水印预览
