package models

import "gorm.io/gorm"

type Filter struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Category string `json:"category" binding:"required"`
	UserID uint `json:"user_id"`
}
