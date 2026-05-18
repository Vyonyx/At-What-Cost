package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vyonyx/at-what-cost/initalisers"
	"github.com/vyonyx/at-what-cost/models"
)

func AddFilter(ctx *gin.Context) {
	var filter models.Filter

	if err := ctx.ShouldBind(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficent filter details"})
		return
	}

	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	filter.UserID = user.ID

	result := initalisers.DB.Create(&filter)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"filter": filter})
}

func GetFilters(ctx *gin.Context) {
	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	var filters []models.Filter
	initalisers.DB.Find(&filters, "user_id = ?", user.ID)

	ctx.JSON(http.StatusOK, filters)
}

type UpdatedFilter struct {
	Name string
	Category string
}

func UpdateFilter(ctx *gin.Context) {
	var body UpdatedFilter
	filterID := ctx.Param("id")

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficent filter details"})
		return
	}

	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	var filter models.Filter

	initalisers.DB.First(&filter, "id = ? AND user_id = ?", filterID, user.ID)

	if filter.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "filter not found under user"})
		return
	}

	filter.Name = body.Name
	filter.Category = body.Category

	result := initalisers.DB.Save(&filter)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, filter)
}

func DeleteFilter(ctx *gin.Context) {
	filterID := ctx.Param("id")

	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	result := initalisers.DB.Delete(&models.Filter{}, "id = ? AND user_id = ?", filterID, user.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete the filter" })
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no filter found to delete"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "filter deleted successfully"})
}
