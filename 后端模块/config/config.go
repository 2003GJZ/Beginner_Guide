package config

import (
	"os"
)

// 数据库配置
type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

// 服务器配置
type ServerConfig struct {
	Port string
}

// 应用配置
type Config struct {
	DB                 DbConfig
	Server             ServerConfig
	JWTKey             string
	CORSAllowedOrigins []string
}

// GetConfig 获取应用配置
func GetConfig() Config {
	// 从环境变量获取配置，如果不存在则使用默认值
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "task_manager")

	serverPort := getEnv("SERVER_PORT", "8080")
	jwtKey := getEnv("JWT_KEY", "your_secret_key_for_jwt_please_change_in_production")

	// 允许的跨域来源
	corsOrigins := []string{"http://localhost:8081"}

	return Config{
		DB: DbConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DbName:   dbName,
		},
		Server: ServerConfig{
			Port: serverPort,
		},
		JWTKey:             jwtKey,
		CORSAllowedOrigins: corsOrigins,
	}
}

// 从环境变量获取值，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return c.DB.User + ":" + c.DB.Password + "@(" + c.DB.Host + ":" + c.DB.Port + ")/" + c.DB.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
}
