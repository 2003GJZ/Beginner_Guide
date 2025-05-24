package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"taskmanager/models"
)

// GetTasks 获取所有任务
func GetTasks(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询参数
	priority := c.Query("priority")
	completed := c.Query("completed")

	// 构建查询
	query := db.Where("user_id = ?", userID)

	// 按优先级筛选
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// 按完成状态筛选
	if completed == "true" {
		query = query.Where("completed = ?", true)
	} else if completed == "false" {
		query = query.Where("completed = ?", false)
	}

	// 获取任务列表，按截止日期和创建时间排序
	var tasks []models.Task
	if err := query.Order("CASE WHEN due_date IS NULL THEN 1 ELSE 0 END, due_date ASC, created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		return
	}

	// 转换为响应模型
	response := make([]models.TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = models.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Completed:   task.Completed,
			Priority:    task.Priority,
			DueDate:     task.DueDate,
			UserID:      task.UserID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// TaskRequest 任务请求模型
// 用于创建和更新任务的请求数据结构
type TaskRequest struct {
	Title       string `json:"title"`       // 任务标题，创建时必填
	Description string `json:"description"` // 任务描述，可选
	Completed   bool   `json:"completed"`   // 是否完成，默认false
	Priority    string `json:"priority"`    // 优先级，可选值为"low", "medium", "high"
	DueDate     string `json:"dueDate"`     // 截止日期，字符串格式，可选
}

// CreateTask 创建新任务
func CreateTask(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 绑定请求数据
	var taskReq TaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务数据"})
		return
	}

	// 创建任务时标题必填
	if taskReq.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务标题不能为空"})
		return
	}

	// 验证优先级
	priority := models.Medium // 默认中等优先级
	if taskReq.Priority != "" {
		switch models.Priority(taskReq.Priority) {
		case models.Low, models.Medium, models.High:
			priority = models.Priority(taskReq.Priority)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的优先级，可选值为: low, medium, high"})
			return
		}
	}

	// 解析截止日期
	var dueDate *time.Time
	if taskReq.DueDate != "" {
		// 尝试多种日期格式
		parsedTime, err := parseTime(taskReq.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日期格式"})
			return
		}
		dueDate = &parsedTime
	}

	// 创建任务模型
	task := models.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
		Completed:   taskReq.Completed,
		Priority:    priority,
		DueDate:     dueDate,
		UserID:      userID.(uint),
	}

	// 保存任务
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	// 转换为响应模型
	response := models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTask 更新任务状态
func UpdateTask(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取任务ID
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 查找任务
	var task models.Task
	if db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限"})
		return
	}

	// 绑定更新数据
	var updateData TaskRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的更新数据"})
		return
	}

	// 更新基本字段
	if updateData.Title != "" {
		task.Title = updateData.Title
	}

	task.Description = updateData.Description
	task.Completed = updateData.Completed

	// 解析截止日期
	if updateData.DueDate != "" {
		parsedTime, err := parseTime(updateData.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日期格式"})
			return
		}
		task.DueDate = &parsedTime
	} else {
		task.DueDate = nil
	}

	// 更新优先级
	if updateData.Priority != "" {
		switch models.Priority(updateData.Priority) {
		case models.Low, models.Medium, models.High:
			task.Priority = models.Priority(updateData.Priority)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的优先级，可选值为: low, medium, high"})
			return
		}
	}

	// 保存更新
	if err := db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务失败"})
		return
	}

	// 转换为响应模型
	response := models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// parseTime 尝试解析多种格式的日期字符串
// 支持的格式包括:
// - 2006-01-02 15:04:05
// - 2006-01-02T15:04:05Z
// - 2006-01-02T15:04:05+08:00
// - 2006-01-02
func parseTime(dateStr string) (time.Time, error) {
	// 先尝试使用time.Parse解析标准格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	var parsedTime time.Time
	var err error
	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateStr)
		if err == nil {
			log.Printf("成功解析日期: %s, 格式: %s", dateStr, format)
			return parsedTime, nil
		}
	}

	// 如果上面的格式都不匹配，尝试使用time.ParseInLocation解析带时区的格式
	parsedTime, err = time.Parse(time.RFC3339, dateStr)
	if err == nil {
		log.Printf("成功解析RFC3339日期: %s", dateStr)
		return parsedTime, nil
	}

	// 尝试解析带毫秒的RFC3339格式
	parsedTime, err = time.Parse(time.RFC3339Nano, dateStr)
	if err == nil {
		log.Printf("成功解析RFC3339Nano日期: %s", dateStr)
		return parsedTime, nil
	}

	// 如果还是失败，尝试去掉时区信息只保留日期部分
	if len(dateStr) >= 10 {
		dateOnly := dateStr[:10]
		parsedTime, err = time.Parse("2006-01-02", dateOnly)
		if err == nil {
			log.Printf("成功解析日期部分: %s -> %s", dateStr, dateOnly)
			return parsedTime, nil
		}
	}

	log.Printf("无法解析日期: %s, 错误: %v", dateStr, err)
	return time.Time{}, err
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取任务ID
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 查找任务
	var task models.Task
	if db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限"})
		return
	}

	// 删除任务
	if err := db.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}
