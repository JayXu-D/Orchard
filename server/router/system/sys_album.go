package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AlbumRouter struct{}

// InitAlbumRouter 初始化相册路由
func (s *AlbumRouter) InitAlbumRouter(Router *gin.RouterGroup) {
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
}
