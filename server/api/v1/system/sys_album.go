package system

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AlbumApi struct{}

// CreateAlbum 创建相册
// @Tags Album
// @Summary 创建相册
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateAlbum true "创建相册"
// @Success 200 {object} response.Response{data=response.AlbumResponse,msg=string} "创建成功"
// @Router /album/create [post]
func (albumApi *AlbumApi) CreateAlbum(c *gin.Context) {
	var albumReq request.CreateAlbum
	err := c.ShouldBindJSON(&albumReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	album, err := albumService.CreateAlbum(albumReq)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}

	albumResponse := systemRes.ToAlbumResponse(album)
	response.OkWithData(albumResponse, c)
}

// DeleteAlbum 删除相册
// @Tags Album
// @Summary 删除相册
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteAlbum true "删除相册"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /album/delete [delete]
func (albumApi *AlbumApi) DeleteAlbum(c *gin.Context) {
	var albumReq request.DeleteAlbum
	err := c.ShouldBindJSON(&albumReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = albumService.DeleteAlbum(albumReq)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// UpdateAlbum 更新相册
// @Tags Album
// @Summary 更新相册
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateAlbum true "更新相册"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /album/update [put]
func (albumApi *AlbumApi) UpdateAlbum(c *gin.Context) {
	var albumReq request.UpdateAlbum
	err := c.ShouldBindJSON(&albumReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = albumService.UpdateAlbum(albumReq)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetAlbum 根据ID获取相册
// @Tags Album
// @Summary 根据ID获取相册
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetAlbumByID true "根据ID获取相册"
// @Success 200 {object} response.Response{data=response.AlbumResponse,msg=string} "获取成功"
// @Router /album/get [post]
func (albumApi *AlbumApi) GetAlbum(c *gin.Context) {
	var albumReq request.GetAlbumByID
	err := c.ShouldBindJSON(&albumReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	album, err := albumService.GetAlbum(albumReq)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	albumResponse := systemRes.ToAlbumResponse(album)
	response.OkWithData(albumResponse, c)
}

// GetAlbumList 获取相册列表
// @Tags Album
// @Summary 获取相册列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetAlbumList true "获取相册列表"
// @Success 200 {object} response.Response{data=response.AlbumListResponse,msg=string} "获取成功"
// @Router /album/list [post]
func (albumApi *AlbumApi) GetAlbumList(c *gin.Context) {
	var pageInfo request.GetAlbumList
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := albumService.GetAlbumList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	albumListResponse := systemRes.ToAlbumListResponse(list, total)
	response.OkWithData(albumListResponse, c)
}

// GetAlbumsByCreator 根据创建者UUID获取相册列表
// @Tags Album
// @Summary 根据创建者UUID获取相册列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param creatorUUID path string true "创建者UUID"
// @Success 200 {object} response.Response{data=[]response.AlbumResponse,msg=string} "获取成功"
// @Router /album/creator/{creatorUUID} [get]
func (albumApi *AlbumApi) GetAlbumsByCreator(c *gin.Context) {
	creatorUUID := c.Param("creatorUUID")
	if creatorUUID == "" {
		response.FailWithMessage("创建者UUID不能为空", c)
		return
	}

	// 解析UUID
	uuid, err := uuid.Parse(creatorUUID)
	if err != nil {
		response.FailWithMessage("无效的UUID格式", c)
		return
	}

	list, err := albumService.GetAlbumsByCreator(uuid)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	albumResponses := make([]systemRes.AlbumResponse, len(list))
	for i, album := range list {
		albumResponses[i] = systemRes.ToAlbumResponse(album)
	}

	response.OkWithData(albumResponses, c)
}

// GetAlbumsByAdmin 根据管理员ID获取相册列表
// @Tags Album
// @Summary 根据管理员ID获取相册列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param adminID path int true "管理员ID"
// @Success 200 {object} response.Response{data=[]response.AlbumResponse,msg=string} "获取成功"
// @Router /album/admin/{adminID} [get]
func (albumApi *AlbumApi) GetAlbumsByAdmin(c *gin.Context) {
	adminIDStr := c.Param("adminID")
	if adminIDStr == "" {
		response.FailWithMessage("管理员ID不能为空", c)
		return
	}

	// 解析管理员ID
	var adminID uint
	_, err := fmt.Sscanf(adminIDStr, "%d", &adminID)
	if err != nil {
		response.FailWithMessage("无效的管理员ID格式", c)
		return
	}

	list, err := albumService.GetAlbumsByAdmin(adminID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	albumResponses := make([]systemRes.AlbumResponse, len(list))
	for i, album := range list {
		albumResponses[i] = systemRes.ToAlbumResponse(album)
	}

	response.OkWithData(albumResponses, c)
}
