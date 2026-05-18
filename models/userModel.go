package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique" binding:"required,email"`
	Password string `json:"-" binding:"required,min=8"`
	Filters []Filter `json:"filters,omitempty"`
	Savings []Saving `json:"savings,omitempty"`
}
