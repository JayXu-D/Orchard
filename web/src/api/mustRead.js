import service from '@/utils/request'

// 创建必读内容
export const createMustRead = (data) => {
  return service({
    url: '/mustRead/create',
    method: 'post',
    data
  })
}

// 更新必读内容
export const updateMustRead = (data) => {
  return service({
    url: '/mustRead/update',
    method: 'put',
    data
  })
}

// 删除必读内容
export const deleteMustRead = (data) => {
  return service({
    url: '/mustRead/delete',
    method: 'delete',
    data
  })
}

// 根据ID获取必读内容
export const getMustRead = (data) => {
  return service({
    url: '/mustRead/get',
    method: 'post',
    data
  })
}

// 获取最新必读内容
export const getLatestMustRead = () => {
  return service({
    url: '/mustRead/latest',
    method: 'get'
  })
}
