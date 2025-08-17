# 相册API测试文档

## 1. 创建相册

**接口地址：** `POST /api/album/create`

**请求参数：**
```json
{
  "creatorUUID": "550e8400-e29b-41d4-a716-446655440000",
  "title": "我的相册",
  "coverImageURL": "https://example.com/cover.jpg",
  "description": "这是一个测试相册",
  "adminUserIDs": [1, 2, 3]
}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "creatorUUID": "550e8400-e29b-41d4-a716-446655440000",
    "title": "我的相册",
    "coverImageURL": "https://example.com/cover.jpg",
    "description": "这是一个测试相册",
    "status": 1,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "creator": {
      "id": 1,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "username": "admin",
      "nickName": "管理员",
      "headerImg": "https://example.com/avatar.jpg"
    },
    "adminUsers": [
      {
        "id": 1,
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "username": "admin",
        "nickName": "管理员",
        "headerImg": "https://example.com/avatar.jpg"
      }
    ]
  },
  "msg": "创建成功"
}
```

## 2. 获取相册列表

**接口地址：** `POST /api/album/list`

**请求参数：**
```json
{
  "page": 1,
  "pageSize": 10,
  "title": "",
  "creatorUUID": "",
  "status": 0
}
```

## 3. 根据ID获取相册

**接口地址：** `POST /api/album/get`

**请求参数：**
```json
{
  "id": 1
}
```

## 4. 更新相册

**接口地址：** `PUT /api/album/update`

**请求参数：**
```json
{
  "id": 1,
  "title": "更新后的相册标题",
  "coverImageURL": "https://example.com/new-cover.jpg",
  "description": "更新后的描述",
  "status": 1,
  "adminUserIDs": [1, 2]
}
```

## 5. 删除相册

**接口地址：** `DELETE /api/album/delete`

**请求参数：**
```json
{
  "id": 1
}
```

## 6. 根据创建者UUID获取相册列表

**接口地址：** `GET /api/album/creator/{creatorUUID}`

**路径参数：**
- creatorUUID: 创建者的UUID

## 7. 根据管理员ID获取相册列表

**接口地址：** `GET /api/album/admin/{adminID}`

**路径参数：**
- adminID: 管理员ID

## 注意事项

1. 所有接口都需要JWT认证
2. creatorUUID必须是有效的用户UUID
3. adminUserIDs数组中的用户ID必须是存在的用户
4. 相册状态：1=正常，2=禁用
5. 创建相册时会自动设置状态为1（正常） 