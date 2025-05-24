package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/crypto/bcrypt"

	"taskmanager/config"
	"taskmanager/models"
)

// 密钥，实际应用中应该从配置文件中读取
var jwtKey = []byte("your_secret_key")

// 定义JWT的Claims结构
type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

// RegisterRequest 用户注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	// 定义请求结构体
	var registerReq RegisterRequest

	// 打印原始请求体以调试
	rawData, _ := io.ReadAll(c.Request.Body)
	log.Printf("原始请求体: %s", string(rawData))

	// 重新设置请求体，因为读取了原始数据
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	// 绑定JSON数据到请求结构体
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		log.Println("注册数据绑定失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	log.Printf("尝试注册用户: %s, 密码长度: %d", registerReq.Username, len(registerReq.Password))

	// 检查密码是否为空
	if len(registerReq.Password) == 0 {
		log.Println("密码为空，请求被拒绝")
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	// 创建用户对象
	user := models.User{
		Username: registerReq.Username,
		Password: registerReq.Password,
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if !db.Where("username = ?", user.Username).First(&existingUser).RecordNotFound() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	log.Printf("开始加密密码为用户: %s, 原始密码长度: %d", user.Username, len(user.Password))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	user.Password = string(hashedPassword)
	log.Printf("密码加密成功, 加密后长度: %d", len(user.Password))

	// 创建用户
	log.Printf("开始创建用户: %s", user.Username)
	if err := db.Create(&user).Error; err != nil {
		log.Printf("用户创建失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	log.Printf("用户注册成功: %s, ID: %d", user.Username, user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println("登录数据绑定失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	log.Printf("尝试登录用户: %s", loginData.Username)

	// 查找用户
	var user models.User
	if db.Where("username = ?", loginData.Username).First(&user).RecordNotFound() {
		log.Printf("用户不存在: %s", loginData.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	log.Printf("用户存在: %s, 开始验证密码", user.Username)

	// 验证密码
	log.Printf("开始验证密码, 存储的密码哈希长度: %d, 输入密码长度: %d", len(user.Password), len(loginData.Password))

	// 如果密码存储有问题，尝试直接比较
	if len(user.Password) < 20 { // bcrypt哈希通常大于60个字符
		log.Printf("存储的密码可能未加密，尝试直接比较")
		if user.Password != loginData.Password {
			log.Printf("直接密码比较失败")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		log.Printf("直接密码比较成功")
	} else {
		// 正常的bcrypt比较
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
		if err != nil {
			log.Printf("密码验证失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		log.Printf("密码验证成功通过bcrypt")
	}

	log.Printf("密码验证成功为用户: %s", user.Username)

	// 创建JWT Token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	log.Printf("开始生成JWT令牌为用户ID: %d", user.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("JWT令牌生成失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成token"})
		return
	}

	log.Printf("JWT令牌生成成功为用户: %s", user.Username)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user models.User
	if db.Where("id = ?", userID).First(&user).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 构建头像URL
	avatarUrl := ""
	if user.AvatarPath != "" {
		// 获取MinIO配置
		minioConfig := config.GetMinioConfig()
		avatarUrl = fmt.Sprintf("http://%s/%s/%s", minioConfig.Endpoint, minioConfig.Bucket, user.AvatarPath)
	}

	// 返回用户信息，不包含敏感信息
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"avatarUrl": avatarUrl,
		"createdAt": user.CreatedAt,
	})
}

// UploadAvatar 上传用户头像
func UploadAvatar(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文件
	file, fileHeader, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取头像文件失败"})
		return
	}
	defer file.Close()

	// 检查文件大小（限制为2MB）
	if fileHeader.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "头像文件大小不能超过2MB"})
		return
	}

	// 生成唯一的文件名
	avatarFileName := fmt.Sprintf("avatar_%d_%s%s", userID, uuid.New().String(), filepath.Ext(fileHeader.Filename))

	// 获取MinIO配置
	minioConfig := config.GetMinioConfig()

	// 上传文件到MinIO
	contentType := "image/jpeg"
	if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".gif") {
		contentType = "image/gif"
	}

	_, err = config.MinioClient.PutObject(
		context.Background(),
		minioConfig.Bucket,
		avatarFileName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传头像失败"})
		return
	}

	// 查询用户
	var user models.User
	if db.Where("id = ?", userID).First(&user).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 如果用户已有头像，删除旧头像
	if user.AvatarPath != "" {
		// 尝试删除旧头像，忽略错误
		config.MinioClient.RemoveObject(
			context.Background(),
			minioConfig.Bucket,
			user.AvatarPath,
			minio.RemoveObjectOptions{},
		)
	}

	// 更新用户头像路径
	user.AvatarPath = avatarFileName
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户头像失败"})
		return
	}

	// 构建头像URL
	avatarUrl := fmt.Sprintf("http://%s/%s/%s", minioConfig.Endpoint, minioConfig.Bucket, avatarFileName)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":   "头像上传成功",
		"avatarUrl": avatarUrl,
	})
}

// GetAvatar 获取用户头像
func GetAvatar(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询用户信息
	var user models.User
	if db.Where("id = ?", userID).First(&user).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 构建头像URL
	avatarUrl := ""
	if user.AvatarPath != "" {
		// 获取MinIO配置
		minioConfig := config.GetMinioConfig()
		avatarUrl = fmt.Sprintf("http://%s/%s/%s", minioConfig.Endpoint, minioConfig.Bucket, user.AvatarPath)
	}

	c.JSON(http.StatusOK, gin.H{"avatarUrl": avatarUrl})
}
