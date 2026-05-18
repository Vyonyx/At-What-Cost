package models

import (
	"time"

	"gorm.io/gorm"
)

type Saving struct {
	gorm.Model
	Name string `json:"name"`
	Description string `json:"description"`
	DepositFrequency string `json:"deposit_frequency"`
	DepositAmount int `json:"deposit_amount"`
	TotalAmount int `json:"total_amount"`
	LastCalculatedAt time.Time `json:"last_calculated_at"`
	UserID uint `json:"user_id"`
}
