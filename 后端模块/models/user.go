package models

import (
	"github.com/jinzhu/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username   string `gorm:"unique;not null" json:"username"`
	Password   string `gorm:"not null" json:"-"` // 密码不会在JSON中返回
	Email      string `gorm:"size:100" json:"email"`
	AvatarPath string `gorm:"size:255" json:"avatar_path"` // 头像存储路径
}
