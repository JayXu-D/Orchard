package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AlbumRouter struct{}

// InitAlbumRouter 初始化相册路由
func (s *AlbumRouter) InitAlbumRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	albumRouter := Router.Group("album").Use(middleware.OperationRecord())
	albumRouterWithoutRecord := Router.Group("album")
	{
		albumRouter.POST("create", albumApi.CreateAlbum)   // 创建相册
		albumRouter.DELETE("delete", albumApi.DeleteAlbum) // 删除相册
		albumRouter.PUT("update", albumApi.UpdateAlbum)    // 更新相册
	}
	{
		albumRouterWithoutRecord.POST("get", albumApi.GetAlbum)                           // 根据ID获取相册
		albumRouterWithoutRecord.POST("list", albumApi.GetAlbumList)                      // 获取相册列表
		albumRouterWithoutRecord.GET("creator/:creatorUUID", albumApi.GetAlbumsByCreator) // 根据创建者UUID获取相册列表
		albumRouterWithoutRecord.GET("admin/:adminID", albumApi.GetAlbumsByAdmin)         // 根据管理员ID获取相册列表
	}

	// 图纸路由
	drawingRouter := Router.Group("drawing").Use(middleware.OperationRecord())
	drawingRouterWithoutRecord := Router.Group("drawing")
	{
		drawingRouter.POST("create", drawingApi.CreateDrawing)   // 创建图纸
		drawingRouter.DELETE("delete", drawingApi.DeleteDrawing) // 删除图纸
		drawingRouter.PUT("update", drawingApi.UpdateDrawing)    // 更新图纸
	}
	{
		drawingRouterWithoutRecord.POST("get", drawingApi.GetDrawingByID)                  // 根据ID获取图纸
		drawingRouterWithoutRecord.POST("list", drawingApi.GetDrawingList)                 // 获取图纸列表
		drawingRouterWithoutRecord.POST("download", drawingApi.DownloadDrawing)            // 下载图纸
		drawingRouterWithoutRecord.POST("batchDownload", drawingApi.BatchDownloadDrawings) // 批量下载图纸
		drawingRouterWithoutRecord.GET("watermark/:filename", drawingApi.GetWatermarkFile) // 获取水印文件
		drawingRouterWithoutRecord.GET("file/:filename", drawingApi.GetDrawingFile)        // 获取图纸文件
	}

	// 在公共路由组中添加v1路径的文件访问路由，避免权限认证问题
	v1DrawingRouter := PublicRouter.Group("v1").Group("drawing")
	{
		v1DrawingRouter.GET("watermark/:filename", drawingApi.GetWatermarkFile) // 获取水印文件
		v1DrawingRouter.GET("file/:filename", drawingApi.GetDrawingFile)        // 获取图纸文件
	}
}
