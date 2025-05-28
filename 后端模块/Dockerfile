FROM golang:latest

WORKDIR /app

# 安装等待脚本所需的工具
RUN apt-get update && apt-get install -y netcat-openbsd

# 复制等待脚本
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN go build -o taskmanager .

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

# 设置环境变量
ENV DB_HOST=mysql \
    DB_PORT=3306 \
    DB_USER=root \
    DB_PASSWORD=root \
    DB_NAME=task_manager \
    SERVER_PORT=8080 \
    JWT_KEY=your_secret_key_for_jwt_please_change_in_production \
    MINIO_ENDPOINT=minio:9000 \
    MINIO_ACCESS_KEY=admin \
    MINIO_SECRET_KEY=12345678 \
    MINIO_BUCKET=taskmanager

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["/app/taskmanager"]
