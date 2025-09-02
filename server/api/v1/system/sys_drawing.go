package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DrawingApi struct{}

// CreateDrawing 创建图纸
// @Tags Drawing
// @Summary 创建图纸
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateDrawing true "创建图纸"
// @Success 200 {object} response.Response{data=response.DrawingResponse,msg=string} "创建成功"
// @Router /drawing/create [post]
func (drawingApi *DrawingApi) CreateDrawing(c *gin.Context) {
	var drawingReq request.CreateDrawing
	err := c.ShouldBindJSON(&drawingReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	drawing, err := drawingService.CreateDrawing(drawingReq)
	if err != nil {
		global.GVA_LOG.Error("创建图纸失败!", zap.Error(err))
		response.FailWithMessage("创建图纸失败", c)
		return
	}

	drawingResponse := systemRes.ToDrawingResponse(drawing)
	response.OkWithData(drawingResponse, c)
}

// UpdateDrawing 更新图纸
// @Tags Drawing
// @Summary 更新图纸
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateDrawing true "更新图纸"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /drawing/update [put]
func (drawingApi *DrawingApi) UpdateDrawing(c *gin.Context) {
	var drawingReq request.UpdateDrawing
	err := c.ShouldBindJSON(&drawingReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = drawingService.UpdateDrawing(drawingReq)
	if err != nil {
		global.GVA_LOG.Error("更新图纸失败!", zap.Error(err))
		response.FailWithMessage("更新图纸失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteDrawing 删除图纸
// @Tags Drawing
// @Summary 删除图纸
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteDrawing true "删除图纸"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /drawing/delete [delete]
func (drawingApi *DrawingApi) DeleteDrawing(c *gin.Context) {
	var drawingReq request.DeleteDrawing
	err := c.ShouldBindJSON(&drawingReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = drawingService.DeleteDrawing(drawingReq)
	if err != nil {
		global.GVA_LOG.Error("删除图纸失败!", zap.Error(err))
		response.FailWithMessage("删除图纸失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetDrawingByID 根据ID获取图纸
// @Tags Drawing
// @Summary 根据ID获取图纸
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetDrawingByID true "根据ID获取图纸"
// @Success 200 {object} response.Response{data=response.DrawingResponse,msg=string} "获取成功"
// @Router /drawing/get [post]
func (drawingApi *DrawingApi) GetDrawingByID(c *gin.Context) {
	var drawingReq request.GetDrawingByID
	err := c.ShouldBindJSON(&drawingReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	drawing, err := drawingService.GetDrawingByID(drawingReq)
	if err != nil {
		global.GVA_LOG.Error("获取图纸失败!", zap.Error(err))
		response.FailWithMessage("获取图纸失败", c)
		return
	}

	drawingResponse := systemRes.ToDrawingResponse(drawing)
	response.OkWithData(drawingResponse, c)
}

// GetDrawingList 获取图纸列表
// @Tags Drawing
// @Summary 获取图纸列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetDrawingList true "获取图纸列表"
// @Success 200 {object} response.Response{data=response.DrawingListResponse,msg=string} "获取成功"
// @Router /drawing/list [post]
func (drawingApi *DrawingApi) GetDrawingList(c *gin.Context) {
	var pageInfo request.GetDrawingList
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := drawingService.GetDrawingList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取图纸列表失败!", zap.Error(err))
		response.FailWithMessage("获取图纸列表失败", c)
		return
	}

	drawingListResponse := systemRes.ToDrawingListResponse(list, total)
	response.OkWithData(drawingListResponse, c)
}

// DownloadDrawing 下载图纸
// @Tags Drawing
// @Summary 下载图纸
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DownloadDrawing true "下载图纸"
// @Success 200 {object} response.Response{data=response.DownloadResponse,msg=string} "下载成功"
// @Router /drawing/download [post]
func (drawingApi *DrawingApi) DownloadDrawing(c *gin.Context) {
	var downloadReq request.DownloadDrawing
	err := c.ShouldBindJSON(&downloadReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从JWT中获取用户UUID
	userUUID := utils.GetUserUuid(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户身份验证失败", c)
		return
	}

	downloadResponse, err := drawingService.DownloadDrawing(downloadReq, userUUID)
	if err != nil {
		global.GVA_LOG.Error("下载图纸失败!", zap.Error(err))
		response.FailWithMessage("下载图纸失败", c)
		return
	}

	response.OkWithData(downloadResponse, c)
}

// BatchDownloadDrawings 批量下载图纸
// @Tags Drawing
// @Summary 批量下载图纸
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BatchDownloadDrawings true "批量下载图纸"
// @Success 200 {object} response.Response{data=response.DownloadResponse,msg=string} "下载成功"
// @Router /drawing/batchDownload [post]
func (drawingApi *DrawingApi) BatchDownloadDrawings(c *gin.Context) {
	var batchDownloadReq request.BatchDownloadDrawings
	err := c.ShouldBindJSON(&batchDownloadReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从JWT中获取用户UUID
	userUUID := utils.GetUserUuid(c)
	if userUUID == uuid.Nil {
		response.FailWithMessage("用户身份验证失败", c)
		return
	}

	downloadResponse, err := drawingService.BatchDownloadDrawings(batchDownloadReq, userUUID)
	if err != nil {
		global.GVA_LOG.Error("批量下载图纸失败!", zap.Error(err))
		response.FailWithMessage("批量下载图纸失败", c)
		return
	}

	response.OkWithData(downloadResponse, c)
}

// GetWatermarkFile 获取水印文件
func (drawingApi *DrawingApi) GetWatermarkFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.FailWithMessage("文件名不能为空", c)
		return
	}

	// 构建水印文件路径 - 使用相对于项目根目录的路径
	// 从当前工作目录开始构建路径
	workDir, err := os.Getwd()
	if err != nil {
		global.GVA_LOG.Error("获取工作目录失败", zap.Error(err))
		response.FailWithMessage("服务器内部错误", c)
		return
	}

	filePath := filepath.Join(workDir, "cache", "watermark", filename)

	global.GVA_LOG.Info("获取水印文件",
		zap.String("filename", filename),
		zap.String("work_dir", workDir),
		zap.String("file_path", filePath))

	// 检查缓存目录是否存在
	cacheDir := filepath.Join(workDir, "cache", "watermark")
	if _, err := os.Stat(cacheDir); err != nil {
		global.GVA_LOG.Warn("水印缓存目录不存在",
			zap.String("cache_dir", cacheDir),
			zap.Error(err))

		// 尝试列出工作目录下的内容
		if entries, err := os.ReadDir(workDir); err == nil {
			var dirs []string
			for _, entry := range entries {
				if entry.IsDir() {
					dirs = append(dirs, entry.Name())
				}
			}
			global.GVA_LOG.Info("工作目录下的目录",
				zap.String("work_dir", workDir),
				zap.Strings("directories", dirs))
		}

		response.FailWithMessage("缓存目录不存在", c)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		global.GVA_LOG.Warn("水印文件不存在",
			zap.String("file_path", filePath),
			zap.Error(err))

		// 尝试列出缓存目录下的文件
		if entries, err := os.ReadDir(cacheDir); err == nil {
			var files []string
			for _, entry := range entries {
				if !entry.IsDir() {
					files = append(files, entry.Name())
				}
			}
			global.GVA_LOG.Info("缓存目录下的文件",
				zap.String("cache_dir", cacheDir),
				zap.Strings("files", files))
		}

		response.FailWithMessage("文件不存在", c)
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "image/jpeg")

	// 返回文件
	c.File(filePath)
}

// GetMyDrawings 获取当前用户可下载的图纸列表
// @Tags Drawing
// @Summary 获取当前用户可下载的图纸列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetMyDrawings true "获取我的图纸列表"
// @Success 200 {object} response.Response{data=response.DrawingListResponse,msg=string} "获取成功"
// @Router /drawing/my [post]
func (drawingApi *DrawingApi) GetMyDrawings(c *gin.Context) {
	var pageInfo request.GetMyDrawings
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从JWT claims中获取当前用户信息
	claims, exists := c.Get("claims")
	if !exists {
		response.FailWithMessage("用户信息获取失败", c)
		return
	}

	userClaims := claims.(*request.CustomClaims)
	pageInfo.UserID = userClaims.BaseClaims.ID
	pageInfo.UserUUID = userClaims.BaseClaims.UUID.String()

	list, total, err := drawingService.GetMyDrawings(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取我的图纸列表失败!", zap.Error(err))
		response.FailWithMessage("获取我的图纸列表失败", c)
		return
	}

	drawingListResponse := systemRes.ToDrawingListResponse(list, total)
	response.OkWithData(drawingListResponse, c)
}

// UpdateEmptyDrawings 更新空白图纸记录（临时方法，用于测试）
// @Tags Drawing
// @Summary 更新空白图纸记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /drawing/updateEmpty [post]
func (drawingApi *DrawingApi) UpdateEmptyDrawings(c *gin.Context) {
	err := drawingService.UpdateEmptyDrawings()
	if err != nil {
		global.GVA_LOG.Error("更新空白图纸失败!", zap.Error(err))
		response.FailWithMessage("更新空白图纸失败", c)
		return
	}

	response.OkWithMessage("空白图纸更新成功", c)
}

// GetDrawingFile 获取图纸文件
// @Tags Drawing
// @Summary 获取图纸文件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param filename path string true "文件名"
// @Success 200 {file} file "文件"
// @Router /drawing/file/{filename} [get]
func (drawingApi *DrawingApi) GetDrawingFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.FailWithMessage("文件名不能为空", c)
		return
	}

	// 构建图纸文件路径 - 使用相对于项目根目录的路径
	// 从当前工作目录开始构建路径
	workDir, err := os.Getwd()
	if err != nil {
		global.GVA_LOG.Error("获取工作目录失败", zap.Error(err))
		response.FailWithMessage("服务器内部错误", c)
		return
	}

	filePath := filepath.Join(workDir, "uploads", "file", filename)

	global.GVA_LOG.Info("获取图纸文件",
		zap.String("filename", filename),
		zap.String("work_dir", workDir),
		zap.String("file_path", filePath))

	// 检查uploads目录是否存在
	uploadsDir := filepath.Join(workDir, "uploads", "file")
	if _, err := os.Stat(uploadsDir); err != nil {
		global.GVA_LOG.Warn("uploads目录不存在",
			zap.String("uploads_dir", uploadsDir),
			zap.Error(err))

		// 尝试列出工作目录下的内容
		if entries, err := os.ReadDir(workDir); err == nil {
			var dirs []string
			for _, entry := range entries {
				if entry.IsDir() {
					dirs = append(dirs, entry.Name())
				}
			}
			global.GVA_LOG.Info("工作目录下的目录",
				zap.String("work_dir", workDir),
				zap.Strings("directories", dirs))
		}

		response.FailWithMessage("uploads目录不存在", c)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		global.GVA_LOG.Warn("图纸文件不存在",
			zap.String("file_path", filePath),
			zap.Error(err))

		// 尝试列出uploads目录下的文件
		if entries, err := os.ReadDir(uploadsDir); err == nil {
			var files []string
			for _, entry := range entries {
				if !entry.IsDir() {
					files = append(files, entry.Name())
				}
			}
			global.GVA_LOG.Info("uploads目录下的文件",
				zap.String("uploads_dir", uploadsDir),
				zap.Strings("files", files))
		}

		response.FailWithMessage("文件不存在", c)
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// 根据文件扩展名设置Content-Type
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		c.Header("Content-Type", "image/jpeg")
	case ".png":
		c.Header("Content-Type", "image/png")
	case ".pdf":
		c.Header("Content-Type", "application/pdf")
	case ".dwg":
		c.Header("Content-Type", "application/acad")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	// 返回文件
	c.File(filePath)
}
