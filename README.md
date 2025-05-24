# 任务管理系统示范项目

这是一个前后端分离的任务管理系统示范项目，包含用户注册登录、任务管理、文件上传和用户头像管理功能。项目具有完整的任务管理功能，包括任务的优先级和截止日期设置，以及文件管理和用户个性化功能，适合新人学习前后端开发。

## 项目结构

```
Beginner_Guide/
├── 前端模块/        # Vue.js前端项目
└── 后端模块/        # Go后端项目
```

## 后端项目

### 技术栈

- 语言：Go
- Web框架：Gin
- ORM：GORM
- 数据库：MySQL
- 认证：JWT
- 文件存储：MinIO

### 目录结构

```
后端模块/
├── api/            # API定义
├── config/         # 配置文件
│   ├── config.go   # 数据库配置
│   └── minio.go    # MinIO配置
├── controllers/    # 控制器
│   ├── fileController.go  # 文件管理控制器
│   ├── taskController.go  # 任务管理控制器
│   └── userController.go  # 用户管理控制器
├── middleware/     # 中间件
├── models/         # 数据模型
├── go.mod          # Go模块文件
└── main.go         # 主程序入口
```

### 运行方法

### 后端运行

1. 确保已安装Go环境、MySQL数据库和MinIO服务
2. 创建数据库：
   ```sql
   CREATE DATABASE task_manager CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
3. 可以通过环境变量或修改配置文件来配置连接信息：
   - 数据库配置（`config/config.go`）：
     - DB_HOST: 数据库主机（默认localhost）
     - DB_PORT: 数据库端口（默认3306）
     - DB_USER: 数据库用户名（默认root）
     - DB_PASSWORD: 数据库密码（默认root）
     - DB_NAME: 数据库名称（默认task_manager）
     - SERVER_PORT: 服务器端口（默认8080）
     - JWT_KEY: JWT密钥（生产环境必须修改）
   - MinIO配置（`config/minio.go`）：
     - MINIO_ENDPOINT: MinIO服务地址（默认localhost:9000）
     - MINIO_ACCESS_KEY: MinIO访问密钥（默认admin）
     - MINIO_SECRET_KEY: MinIO密钥（默认12345678）
     - MINIO_BUCKET: MinIO存储桶名称（默认taskfiles）
4. 在后端项目根目录下执行：
   ```bash
   go mod tidy
   go run main.go
   ```
5. 后端服务器将在`http://localhost:8080`上运行

## 前端项目

### 技术栈

- 框架：Vue.js 2
- UI组件库：Element UI
- 状态管理：Vuex
- 路由：Vue Router
- HTTP客户端：Axios

### 目录结构

```
前端模块/
├── public/         # 静态资源
├── src/            # 源代码
│   ├── components/ # 组件
│   │   └── UserAvatar.vue # 用户头像组件
│   ├── router/     # 路由配置
│   ├── store/      # 状态管理
│   ├── views/      # 视图组件
│   │   ├── FileManager.vue # 文件管理页面
│   │   ├── Home.vue        # 主页/任务管理
│   │   ├── Login.vue       # 登录页面
│   │   ├── Profile.vue     # 用户资料页面
│   │   └── Register.vue    # 注册页面
│   ├── App.vue     # 根组件
│   └── main.js     # 入口文件
└── package.json    # 项目配置
```

### 前端运行

1. 确保已安装Node.js和npm
2. 在前端项目根目录下执行：
   ```bash
   npm install
   npm run serve
   ```
3. 前端应用将在`http://localhost:8081`上运行，并通过代理自动连接到后端服务

## API接口说明

### 用户相关

- `POST /api/register`：用户注册
- `POST /api/login`：用户登录
- `GET /api/user/info`：获取用户信息（需要认证）
- `POST /api/user/avatar`：上传用户头像（需要认证）

### 任务相关

- `GET /api/tasks`：获取任务列表（需要认证）
- `POST /api/task`：创建新任务（需要认证）
- `PUT /api/task/:id`：更新任务（需要认证）
- `DELETE /api/task/:id`：删除任务（需要认证）

### 文件管理

- `POST /api/file/upload`：上传文件（需要认证）
- `GET /api/files`：获取文件列表（需要认证）
- `GET /api/file/:filename`：下载文件（需要认证）
- `POST /api/file/delete/:filename`：删除文件（需要认证）

## 系统功能

1. 用户管理：
   - 注册、登录和用户信息获取
   - 用户头像上传和显示
   - 个人资料页面
2. 任务管理：
   - 创建、查看、编辑和删除任务
   - 任务属性：
     - 任务标题和描述
     - 任务优先级（低、中、高）
     - 任务截止日期
     - 任务完成状态
   - 任务过滤：按状态和优先级过滤任务
3. 文件管理：
   - 文件上传（支持多种文件类型，限制大小10MB）
   - 文件列表查看
   - 文件下载
   - 文件删除
   - 图片和PDF预览

## 注意事项

1. 前端项目运行在8081端口，后端运行在8080端口
2. MinIO服务默认运行在9000端口，控制台在9001端口
3. 文件上传大小限制：普通文件10MB，头像2MB
4. 本项目仅作为示范，生产环境使用需要进一步完善安全措施
5. JWT密钥应在生产环境中妥善保管，可通过环境变量配置
6. 数据库和MinIO配置信息已通过环境变量方式支持，更安全便捷
7. 首次使用需要确保MinIO服务已启动并创建了对应的存储桶
