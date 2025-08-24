package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MustReadApi struct{}

// CreateMustRead 创建必读内容
// @Tags MustRead
// @Summary 创建必读内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateMustRead true "创建必读内容"
// @Success 200 {object} response.Response{data=response.MustReadResponse,msg=string} "创建成功"
// @Router /mustRead/create [post]
func (mustReadApi *MustReadApi) CreateMustRead(c *gin.Context) {
	var mustReadReq request.CreateMustRead
	err := c.ShouldBindJSON(&mustReadReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	mustRead, err := mustReadService.CreateMustRead(mustReadReq)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}

	mustReadResponse := systemRes.ToMustReadResponse(mustRead)
	response.OkWithData(mustReadResponse, c)
}

// DeleteMustRead 删除必读内容
// @Tags MustRead
// @Summary 删除必读内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteMustRead true "删除必读内容"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /mustRead/delete [delete]
func (mustReadApi *MustReadApi) DeleteMustRead(c *gin.Context) {
	var mustReadReq request.DeleteMustRead
	err := c.ShouldBindJSON(&mustReadReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = mustReadService.DeleteMustRead(mustReadReq)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// UpdateMustRead 更新必读内容
// @Tags MustRead
// @Summary 更新必读内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateMustRead true "更新必读内容"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /mustRead/update [put]
func (mustReadApi *MustReadApi) UpdateMustRead(c *gin.Context) {
	var mustReadReq request.UpdateMustRead
	err := c.ShouldBindJSON(&mustReadReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = mustReadService.UpdateMustRead(mustReadReq)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetMustRead 根据ID获取必读内容
// @Tags MustRead
// @Summary 根据ID获取必读内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetMustReadByID true "根据ID获取必读内容"
// @Success 200 {object} response.Response{data=response.MustReadResponse,msg=string} "获取成功"
// @Router /mustRead/get [post]
func (mustReadApi *MustReadApi) GetMustRead(c *gin.Context) {
	var mustReadReq request.GetMustReadByID
	err := c.ShouldBindJSON(&mustReadReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	mustRead, err := mustReadService.GetMustRead(mustReadReq)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	mustReadResponse := systemRes.ToMustReadResponse(mustRead)
	response.OkWithData(mustReadResponse, c)
}

// GetLatestMustRead 获取最新必读内容
// @Tags MustRead
// @Summary 获取最新必读内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=response.MustReadResponse,msg=string} "获取成功"
// @Router /mustRead/latest [get]
func (mustReadApi *MustReadApi) GetLatestMustRead(c *gin.Context) {
	mustRead, err := mustReadService.GetLatestMustRead()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	mustReadResponse := systemRes.ToMustReadResponse(mustRead)
	response.OkWithData(mustReadResponse, c)
}
