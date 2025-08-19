package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/gin-gonic/gin"
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
