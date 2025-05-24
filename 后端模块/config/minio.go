package config

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient 全局Minio客户端
var MinioClient *minio.Client

// MinioConfig MinIO配置
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

// GetMinioConfig 获取MinIO配置
func GetMinioConfig() MinioConfig {
	// 使用config包中的getEnv函数
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "admin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "12345678"
	}

	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "taskmanager"
	}

	return MinioConfig{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    false,
		Bucket:    bucket,
	}
}

// InitMinio 初始化MinIO客户端
func InitMinio() {
	config := GetMinioConfig()

	// 初始化MinIO客户端
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Fatalf("初始化MinIO客户端失败: %v", err)
	}

	// 检查存储桶是否存在，不存在则创建
	exists, err := minioClient.BucketExists(context.Background(), config.Bucket)
	if err != nil {
		log.Fatalf("检查存储桶失败: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("创建存储桶失败: %v", err)
		}
		log.Printf("成功创建存储桶: %s", config.Bucket)

		// 设置存储桶策略，允许公共读取
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + config.Bucket + `/*"]
				}
			]
		}`

		err = minioClient.SetBucketPolicy(context.Background(), config.Bucket, policy)
		if err != nil {
			log.Fatalf("设置存储桶策略失败: %v", err)
		}
	}

	MinioClient = minioClient
	log.Println("MinIO客户端初始化成功")
}
