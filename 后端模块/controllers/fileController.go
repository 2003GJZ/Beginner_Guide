package controllers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"taskmanager/config"
)

// FileResponse 文件响应结构
type FileResponse struct {
	FileName string `json:"fileName"` // 原始文件名
	FileURL  string `json:"fileUrl"`  // 文件访问URL
	FileSize int64  `json:"fileSize"` // 文件大小（字节）
	FileType string `json:"fileType"` // 文件类型
	UploadAt string `json:"uploadAt"` // 上传时间
}

// UploadFile 上传文件处理函数
func UploadFile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文件
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// 检查文件大小（限制为10MB）
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过10MB"})
		return
	}

	// 获取文件扩展名
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))

	// 检查文件类型（可以根据需要添加更多允许的类型）
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".doc": true, ".docx": true,
		".xls": true, ".xlsx": true, ".txt": true,
	}

	if !allowedExts[fileExt] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// 生成唯一的文件名
	fileName := fmt.Sprintf("%d_%s%s", userID, uuid.New().String(), fileExt)

	// 获取MinIO配置
	minioConfig := config.GetMinioConfig()

	// 上传文件到MinIO
	contentType := getContentType(fileExt)
	_, err = config.MinioClient.PutObject(
		context.Background(),
		minioConfig.Bucket,
		fileName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败"})
		return
	}

	// 构建文件URL
	fileURL := fmt.Sprintf("http://%s/%s/%s", minioConfig.Endpoint, minioConfig.Bucket, fileName)

	// 返回文件信息
	response := FileResponse{
		FileName: fileHeader.Filename,
		FileURL:  fileURL,
		FileSize: fileHeader.Size,
		FileType: contentType,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// GetFileList 获取文件列表
func GetFileList(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取MinIO配置
	minioConfig := config.GetMinioConfig()

	// 创建一个前缀，只列出当前用户的文件
	prefix := fmt.Sprintf("%d_", userID)

	// 列出对象
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := config.MinioClient.ListObjects(ctx, minioConfig.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var fileList []FileResponse
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件列表失败"})
			return
		}

		// 提取原始文件名（去掉UUID部分）
		originalName := getOriginalFileName(object.Key)

		// 构建文件URL
		fileURL := fmt.Sprintf("http://%s/%s/%s", minioConfig.Endpoint, minioConfig.Bucket, object.Key)

		fileList = append(fileList, FileResponse{
			FileName: originalName,
			FileURL:  fileURL,
			FileSize: object.Size,
			FileType: getContentType(filepath.Ext(object.Key)),
			UploadAt: object.LastModified.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": fileList})
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文件名
	fileName := c.Param("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	// 检查文件名是否以用户ID开头，确保用户只能删除自己的文件
	prefix := fmt.Sprintf("%d_", userID)
	if !strings.HasPrefix(fileName, prefix) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文件"})
		return
	}

	// 获取MinIO配置
	minioConfig := config.GetMinioConfig()

	// 删除文件
	err := config.MinioClient.RemoveObject(
		context.Background(),
		minioConfig.Bucket,
		fileName,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}

// 根据文件扩展名获取内容类型
func getContentType(fileExt string) string {
	switch strings.ToLower(fileExt) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".doc", ".docx":
		return "application/msword"
	case ".xls", ".xlsx":
		return "application/vnd.ms-excel"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

// 从存储的文件名中提取原始文件名
func getOriginalFileName(key string) string {
	// 文件名格式为: userID_UUID.ext
	// 由于我们没有存储原始文件名，这里只返回UUID部分
	parts := strings.Split(key, "_")
	if len(parts) < 2 {
		return key
	}
	return parts[1]
}
