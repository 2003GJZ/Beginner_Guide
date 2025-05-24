package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Priority 任务优先级
type Priority string

// 任务优先级枚举
const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

// Task 任务模型
type Task struct {
	gorm.Model
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Completed   bool       `gorm:"default:false" json:"completed"`
	Priority    Priority   `gorm:"type:varchar(10);default:'medium'" json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
	UserID      uint       `json:"userId"` // 关联到用户
}

// TaskResponse 任务响应模型
type TaskResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
	UserID      uint       `json:"userId"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
