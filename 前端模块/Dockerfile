FROM node:latest

WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm install
RUN npm install -g serve

# 复制项目文件
COPY . .

# 构建项目
RUN npm run build

# 暴露端口
EXPOSE 8081

# 使用serve启动静态文件服务器
CMD ["serve", "-s", "dist", "-l", "8081"]
