package controllers

import (
	"github.com/jinzhu/gorm"
)

// 全局数据库连接
var db *gorm.DB

// SetDB 设置控制器包的数据库连接
func SetDB(database *gorm.DB) {
	db = database
}
