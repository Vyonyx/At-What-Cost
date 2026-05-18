package initalisers

import "github.com/vyonyx/at-what-cost/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Filter{})
	DB.AutoMigrate(&models.Saving{})
}
