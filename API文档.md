# 任务管理系统 API 文档

## 基本信息

- 基础URL: `http://localhost:8080`
- 认证方式: JWT Token（Bearer Authentication）
- 数据格式: JSON
- 请求方式: 按照规范，只使用GET和POST请求
  - GET: 用于获取数据，不需要传递请求体
  - POST: 用于创建、更新和删除数据，需要传递请求体

## 认证

除了登录和注册接口外，所有接口都需要在请求头中包含有效的JWT令牌：

```
Authorization: Bearer <token>
```

## 1. 用户相关接口

### 1.1 用户注册

- **URL**: `/api/register`
- **方法**: `POST`
- **描述**: 创建新用户账号
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- **成功响应** (200):
  ```json
  {
    "message": "用户注册成功"
  }
  ```
- **错误响应**:
  - 400: 请求数据无效或用户名已存在
  - 500: 服务器内部错误

### 1.2 用户登录

- **URL**: `/api/login`
- **方法**: `POST`
- **描述**: 验证用户凭据并返回访问令牌
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- **成功响应** (200):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "用户名"
    }
  }
  ```
- **错误响应**:
  - 400: 请求数据无效
  - 401: 用户名或密码错误
  - 500: 服务器内部错误

### 1.3 获取用户信息

- **URL**: `/api/user/info`
- **方法**: `GET`
- **描述**: 获取当前登录用户的详细信息
- **请求头**: 需要Authorization
- **成功响应** (200):
  ```json
  {
    "id": 1,
    "username": "用户名",
    "email": "user@example.com",
    "avatarUrl": "http://localhost:9000/taskmanager/avatar_1_abc123.jpg",
    "createdAt": "2025-05-24T01:00:00Z"
  }
  ```
- **错误响应**:
  - 401: 未授权
  - 404: 用户不存在
  - 500: 服务器内部错误

### 1.4 上传用户头像

- **URL**: `/api/user/avatar`
- **方法**: `POST`
- **描述**: 上传用户头像图片
- **请求头**: 
  - 需要Authorization
  - Content-Type: multipart/form-data
- **请求参数**:
  - `avatar`: 文件字段，包含要上传的头像图片（支持jpg、png、gif格式，大小不超过2MB）
- **成功响应** (200):
  ```json
  {
    "message": "头像上传成功",
    "avatarUrl": "http://localhost:9000/taskmanager/avatar_1_abc123.jpg"
  }
  ```
- **错误响应**:
  - 400: 无效的文件或文件过大
  - 401: 未授权
  - 500: 服务器内部错误

## 2. 任务相关接口

### 2.1 获取任务列表

- **URL**: `/api/tasks`
- **方法**: `GET`
- **描述**: 获取当前用户的所有任务
- **请求头**: 需要Authorization
- **查询参数**:
  - `priority`: 可选，按优先级筛选（low, medium, high）
  - `completed`: 可选，按完成状态筛选（true, false）
- **成功响应** (200):
  ```json
  [
    {
      "id": 1,
      "title": "任务标题",
      "description": "任务描述",
      "completed": false,
      "priority": "medium",
      "dueDate": "2025-06-01T12:00:00Z",
      "userId": 1,
      "createdAt": "2025-05-24T01:00:00Z",
      "updatedAt": "2025-05-24T01:00:00Z"
    },
    ...
  ]
  ```
- **错误响应**:
  - 401: 未授权
  - 500: 服务器内部错误

### 2.2 创建任务

- **URL**: `/api/task`
- **方法**: `POST`
- **描述**: 创建新任务
- **请求头**: 需要Authorization
- **请求体**:
  ```json
  {
    "title": "任务标题",
    "description": "任务描述",
    "priority": "medium",
    "dueDate": "2025-06-01T12:00:00Z",
    "completed": false
  }
  ```
- **参数说明**:
  - `title`: 必填，任务标题
  - `description`: 可选，任务描述
  - `priority`: 可选，任务优先级，可选值为 "low", "medium", "high"，默认为 "medium"
  - `dueDate`: 可选，任务截止日期，ISO 8601格式
  - `completed`: 可选，任务是否完成，默认为 false
- **成功响应** (200):
  ```json
  {
    "id": 1,
    "title": "任务标题",
    "description": "任务描述",
    "priority": "medium",
    "dueDate": "2025-06-01T12:00:00Z",
    "completed": false,
    "userId": 1,
    "createdAt": "2025-05-24T01:00:00Z",
    "updatedAt": "2025-05-24T01:00:00Z"
  }
  ```
- **错误响应**:
  - 400: 请求数据无效
  - 401: 未授权
  - 500: 服务器内部错误

### 2.3 更新任务

- **URL**: `/api/task/update/{id}`
- **方法**: `POST`
- **描述**: 更新指定ID的任务
- **请求头**: 需要Authorization
- **URL参数**:
  - `id`: 任务ID
- **请求体**:
  ```json
  {
    "title": "更新后的标题",
    "description": "更新后的描述",
    "priority": "high",
    "dueDate": "2025-06-05T18:00:00Z",
    "completed": true
  }
  ```
- **参数说明**:
  - `title`: 可选，任务标题
  - `description`: 可选，任务描述
  - `priority`: 可选，任务优先级，可选值为 "low", "medium", "high"
  - `dueDate`: 可选，任务截止日期，ISO 8601格式
  - `completed`: 可选，任务是否完成
- **成功响应** (200):
  ```json
  {
    "id": 1,
    "title": "更新后的标题",
    "description": "更新后的描述",
    "priority": "high",
    "dueDate": "2025-06-05T18:00:00Z",
    "completed": true,
    "userId": 1,
    "createdAt": "2025-05-24T01:00:00Z",
    "updatedAt": "2025-05-24T02:00:00Z"
  }
  ```
- **错误响应**:
  - 400: 请求数据无效或任务ID无效
  - 401: 未授权
  - 404: 任务不存在或无权限
  - 500: 服务器内部错误

### 2.4 删除任务

- **URL**: `/api/task/delete/{id}`
- **方法**: `POST`
- **描述**: 删除指定ID的任务
- **请求头**: 需要Authorization
- **URL参数**:
  - `id`: 任务ID
- **请求体**: 无需要的请求体，可以发送空对象 `{}`
- **成功响应** (200):
  ```json
  {
    "message": "任务已删除"
  }
  ```
- **错误响应**:
  - 400: 任务ID无效
  - 401: 未授权
  - 404: 任务不存在或无权限
  - 500: 服务器内部错误

## 3. 文件相关接口

### 3.1 上传文件

- **URL**: `/api/file/upload`
- **方法**: `POST`
- **描述**: 上传文件到服务器
- **请求头**: 
  - 需要Authorization
  - Content-Type: multipart/form-data
- **请求参数**:
  - `file`: 文件字段，包含要上传的文件（支持jpg、png、pdf、doc等格式，大小不超过10MB）
- **成功响应** (200):
  ```json
  {
    "fileName": "example.pdf",
    "fileUrl": "http://localhost:9000/taskmanager/1_abc123.pdf",
    "fileSize": 1024,
    "fileType": "application/pdf",
    "uploadAt": "2025-05-24 20:30:45"
  }
  ```
- **错误响应**:
  - 400: 无效的文件或文件过大
  - 401: 未授权
  - 500: 服务器内部错误

### 3.2 获取文件列表

- **URL**: `/api/files`
- **方法**: `GET`
- **描述**: 获取当前用户上传的所有文件
- **请求头**: 需要Authorization
- **成功响应** (200):
  ```json
  {
    "files": [
      {
        "fileName": "example1.pdf",
        "fileUrl": "http://localhost:9000/taskmanager/1_abc123.pdf",
        "fileSize": 1024,
        "fileType": "application/pdf",
        "uploadAt": "2025-05-24 20:30:45"
      },
      {
        "fileName": "example2.jpg",
        "fileUrl": "http://localhost:9000/taskmanager/1_def456.jpg",
        "fileSize": 512,
        "fileType": "image/jpeg",
        "uploadAt": "2025-05-24 20:35:12"
      }
    ]
  }
  ```
- **错误响应**:
  - 401: 未授权
  - 500: 服务器内部错误

### 3.3 删除文件

- **URL**: `/api/file/delete/{fileName}`
- **方法**: `POST`
- **描述**: 删除指定文件
- **请求头**: 需要Authorization
- **路径参数**:
  - `fileName`: 要删除的文件名（存储在服务器上的文件名，不是原始文件名）
- **请求体**:
  ```json
  {}
  ```
  （可以发送空对象）
- **成功响应** (200):
  ```json
  {
    "message": "文件删除成功"
  }
  ```
- **错误响应**:
  - 400: 文件名不能为空
  - 401: 未授权
  - 403: 无权删除此文件
  - 500: 服务器内部错误

## 4. 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "error": "错误描述信息"
}
```

## 4. 注意事项

1. 所有需要认证的接口必须在请求头中包含有效的JWT令牌
2. 任务相关接口只能操作当前用户自己的任务
3. 日期时间格式遵循ISO 8601标准
4. 所有请求和响应的Content-Type均为application/json
5. 按照规范，只使用GET和POST请求，其中GET用于获取数据，POST用于创建、更新和删除数据
6. 前端运行在8081端口，后端运行在8080端口，通过代理进行通信
