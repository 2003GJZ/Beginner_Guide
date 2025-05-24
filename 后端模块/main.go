package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"taskmanager/config"
	"taskmanager/controllers"
	"taskmanager/middleware"
	"taskmanager/models"
)

// 全局数据库连接
var db *gorm.DB

// 全局配置
var appConfig config.Config

func main() {
	// 加载配置
	appConfig = config.GetConfig()

	// 初始化数据库连接
	initDB()
	defer db.Close()

	// 初始化MinIO客户端
	config.InitMinio()

	// 初始化路由
	router := initRouter()

	// 优雅关闭服务器的通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	serverPort := appConfig.Server.Port
	log.Printf("服务器启动在 http://localhost:%s\n", serverPort)
	go func() {
		if err := router.Run(":" + serverPort); err != nil {
			log.Fatalf("服务器启动失败: %v\n", err)
		}
	}()

	// 等待中断信号
	<-quit
	log.Println("正在关闭服务器...")

	// 给服务器5秒的时间完成当前请求
	time.Sleep(5 * time.Second)

	log.Println("服务器已关闭")
}

// 初始化数据库连接
func initDB() {
	var err error
	// 连接MySQL数据库
	dsn := appConfig.GetDSN()
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	// 启用日志
	db.LogMode(true)

	// 自动迁移模式
	db.AutoMigrate(&models.User{}, &models.Task{})

	// 将数据库连接传递给控制器
	controllers.SetDB(db)

	log.Println("数据库连接成功")
}

// 初始化路由
func initRouter() *gin.Engine {
	router := gin.Default()

	// 使用CORS中间件
	router.Use(middleware.CORSMiddleware())

	// 注册路由
	api := router.Group("/api")
	{
		// 用户相关路由
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/user/info", controllers.GetUserInfo)

			// 任务相关路由
			// 按照规范，只使用GET和POST请求
			auth.GET("/tasks", controllers.GetTasks)
			auth.POST("/task", controllers.CreateTask)
			auth.POST("/task/update/:id", controllers.UpdateTask) // 使用POST替代PUT
			auth.POST("/task/delete/:id", controllers.DeleteTask) // 使用POST替代DELETE

			// 文件相关路由
			auth.POST("/file/upload", controllers.UploadFile)           // 上传文件
			auth.GET("/files", controllers.GetFileList)                 // 获取文件列表
			auth.POST("/file/delete/:fileName", controllers.DeleteFile) // 删除文件

			// 头像相关路由
			auth.POST("/user/avatar", controllers.UploadAvatar) // 上传用户头像
		}
	}

	return router
}
