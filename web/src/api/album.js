import service from '@/utils/request'

// 创建相册
export const createAlbum = (data) => {
  return service({
    url: '/album/create',
    method: 'post',
    data
  })
}

// 获取相册列表（分页）
export const getAlbumList = (data) => {
  return service({
    url: '/album/list',
    method: 'post',
    data
  })
}

// 根据创建者UUID获取相册列表（非必要，备用）
export const getAlbumsByCreator = (creatorUUID) => {
  return service({
    url: `/album/creator/${creatorUUID}`,
    method: 'get'
  })
}

// 根据管理员ID获取相册列表（备用）
export const getAlbumsByAdmin = (adminID) => {
  return service({
    url: `/album/admin/${adminID}`,
    method: 'get'
  })
}

// 更新相册
export const updateAlbum = (data) => {
  return service({
    url: '/album/update',
    method: 'put',
    data
  })
}

// 删除相册
export const deleteAlbumApi = (data) => {
  return service({
    url: '/album/delete',
    method: 'delete',
    data
  })
}

// 获取相册详情
export const getAlbumDetail = (data) => {
  return service({
    url: '/album/get',
    method: 'post',
    data
  })
}

// 获取图纸列表
export const getDrawingList = (data) => {
  return service({
    url: '/drawing/list',
    method: 'post',
    data
  })
}

// 下载图纸
export const downloadDrawing = (data) => {
  return service({
    url: '/drawing/download',
    method: 'get',
    data
  })
}

// 创建图纸
export const createDrawing = (data) => {
  return service({
    url: '/drawing/create',
    method: 'post',
    data
  })
}

// 获取图纸详情
export const getDrawingDetail = (data) => {
  return service({
    url: '/drawing/get',
    method: 'post',
    data
  })
}

// 更新图纸
export const updateDrawing = (data) => {
  return service({
    url: '/drawing/update',
    method: 'put',
    data
  })
}

// 删除图纸
export const deleteDrawing = (data) => {
  return service({
    url: '/drawing/delete',
    method: 'delete',
    data
  })
}