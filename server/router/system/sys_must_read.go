package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MustReadRouter struct{}

// InitMustReadRouter 初始化必读内容路由
func (s *MustReadRouter) InitMustReadRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	mustReadRouter := Router.Group("mustRead").Use(middleware.OperationRecord())
	mustReadRouterWithoutRecord := Router.Group("mustRead")
	{
		mustReadRouter.POST("create", mustReadApi.CreateMustRead)   // 创建必读内容
		mustReadRouter.DELETE("delete", mustReadApi.DeleteMustRead) // 删除必读内容
		mustReadRouter.PUT("update", mustReadApi.UpdateMustRead)    // 更新必读内容
	}
	{
		mustReadRouterWithoutRecord.POST("get", mustReadApi.GetMustRead)         // 根据ID获取必读内容
		mustReadRouterWithoutRecord.GET("latest", mustReadApi.GetLatestMustRead) // 获取最新必读内容
	}
}
