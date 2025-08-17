# 相册功能说明

## 功能概述

相册功能允许用户创建、管理相册，并支持多管理员协作。每个相册都有一个创建者和多个管理员。

## 数据库表结构

### 1. sys_albums (相册表)
- `id`: 主键ID
- `creator_uuid`: 创建者UUID (关联sys_users.uuid)
- `title`: 相册标题
- `cover_image_url`: 相册封面图URL
- `description`: 相册描述
- `status`: 相册状态 (1=正常, 2=禁用)
- `created_at`: 创建时间
- `updated_at`: 更新时间

### 2. sys_album_admin (相册管理员关联表)
- `album_id`: 相册ID (主键)
- `user_id`: 用户ID (主键)

## API接口

### 创建相册
- **POST** `/api/album/create`
- 需要参数：creatorUUID, title, coverImageURL, description, adminUserIDs

### 获取相册列表
- **POST** `/api/album/list`
- 支持分页和条件查询

### 根据ID获取相册
- **POST** `/api/album/get`
- 返回完整的相册信息，包括创建者和管理员信息

### 更新相册
- **PUT** `/api/album/update`
- 可以更新标题、封面、描述、状态和管理员

### 删除相册
- **DELETE** `/api/album/delete`
- 删除相册及其管理员关联

### 根据创建者获取相册
- **GET** `/api/album/creator/{creatorUUID}`
- 获取指定用户创建的所有相册

### 根据管理员获取相册
- **GET** `/api/album/admin/{adminID}`
- 获取指定用户作为管理员的所有相册

## 使用示例

### 1. 创建相册
```bash
curl -X POST http://localhost:8888/api/album/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "creatorUUID": "550e8400-e29b-41d4-a716-446655440000",
    "title": "我的相册",
    "coverImageURL": "https://example.com/cover.jpg",
    "description": "这是一个测试相册",
    "adminUserIDs": [1, 2, 3]
  }'
```

### 2. 获取相册列表
```bash
curl -X POST http://localhost:8888/api/album/list \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "page": 1,
    "pageSize": 10
  }'
```

## 权限控制

- 所有接口都需要JWT认证
- 创建相册时需要验证创建者UUID的有效性
- 管理员ID必须是存在的用户ID
- 支持事务操作，确保数据一致性

## 注意事项

1. 创建相册时会自动设置状态为1（正常）
2. 删除相册时会同时删除相关的管理员关联
3. 更新相册时如果提供了adminUserIDs，会完全替换原有的管理员列表
4. 所有UUID必须是有效的格式
5. 支持软删除（通过status字段控制）

## 扩展功能

可以考虑添加以下功能：
- 相册图片管理
- 相册分享功能
- 相册访问权限控制
- 相册评论功能
- 相册标签功能 