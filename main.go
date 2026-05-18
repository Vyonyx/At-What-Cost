package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vyonyx/at-what-cost/controllers"
	"github.com/vyonyx/at-what-cost/initalisers"
	"github.com/vyonyx/at-what-cost/middleware"
)

func init() {
	initalisers.LoadEnvVariables()
	initalisers.ConnectToDb()
	initalisers.SyncDatabase()
}

func main() {
	r := gin.Default()

	auth := r.Group("/", middleware.RequireAuth)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	auth.GET("/validate", controllers.Validate)
	auth.POST("/logout", controllers.Logout)

	filters := r.Group("/filters", middleware.RequireAuth)
	filters.POST("/", controllers.AddFilter)
	filters.GET("/", controllers.GetFilters)
	filters.PUT("/:id", controllers.UpdateFilter)
	filters.DELETE("/:id", controllers.DeleteFilter)

	savings := r.Group("/savings", middleware.RequireAuth)
	savings.POST("/", controllers.AddSaving)
	savings.GET("/", controllers.GetSavings)
	savings.DELETE("/:id", controllers.DeleteSaving)

	r.Run()
}
