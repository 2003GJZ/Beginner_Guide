@echo off
echo 任务管理系统离线设置脚本 - Windows版本
echo ======================================

REM 检查Docker是否安装
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker未安装，请先安装Docker
    exit /b
)

REM 检查Docker是否在运行
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker未运行，请启动Docker Desktop然后重新运行此脚本
    start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    exit /b
)

echo Docker已安装并正在运行

REM 检查MySQL容器是否已运行
docker ps | findstr mysql >nul
if %errorlevel% neq 0 (
    echo MySQL容器未运行，正在启动MySQL容器...
    docker start d802f9e5c9d4689bbe46f074c772124483cbdeb365f84de064d0f4db8b1f813e
    if %errorlevel% neq 0 (
        echo 无法启动MySQL容器，请检查容器ID是否正确
        exit /b
    )
)

REM 检查MinIO容器是否已运行
docker ps | findstr minio >nul
if %errorlevel% neq 0 (
    echo MinIO容器未运行，正在启动MinIO容器...
    docker start fea74abf87000b38fa39d63df0dc8238603676e9f70669b8b648ad2fc4024cf5
    if %errorlevel% neq 0 (
        echo 无法启动MinIO容器，请检查容器ID是否正确
        exit /b
    )
)

echo MySQL和MinIO容器已启动

REM 创建数据库（如果不存在）
echo 正在创建数据库（如果不存在）...
docker exec -i mysql mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS task_manager CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

REM 初始化MinIO
echo 正在初始化MinIO...
docker exec -i minio mkdir -p /data/taskmanager

REM 编译后端
echo 正在编译后端...
cd /d "%~dp0后端模块"
go build -o taskmanager.exe

REM 构建前端
echo 正在构建前端...
cd /d "%~dp0前端模块"
call npm install
call npm run build

echo ======================================
echo 离线设置完成!
echo 您现在可以使用以下命令启动应用:
echo docker-compose -f docker-offline.yml up -d
echo ======================================

pause
