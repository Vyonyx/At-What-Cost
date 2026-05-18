package controllers

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vyonyx/at-what-cost/initalisers"
	"github.com/vyonyx/at-what-cost/models"
)

func AddSaving(ctx *gin.Context) {
	var saving models.Saving

	if err := ctx.ShouldBind(&saving); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient saving details"})
		return
	}

	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	saving.UserID = user.ID

	if saving.LastCalculatedAt.IsZero() {
		saving.LastCalculatedAt = time.Now()
	}

	result := initalisers.DB.Create(&saving)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, saving)
}

func GetSavings(ctx *gin.Context) {
	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	var savings []models.Saving
	initalisers.DB.Find(&savings, "user_id = ?", user.ID)

	if len(savings) > 0 {
		lazyUpdateSavings(&savings)
	}

	ctx.JSON(http.StatusOK, savings)
}

func UpdateSaving(ctx *gin.Context) {
	var body models.Saving
	savingID := ctx.Param("id")

	if ctx.ShouldBind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient saving details"})
		return
	}

	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	var saving models.Saving

	initalisers.DB.First(&saving, "id = ? AND user_id = ?", savingID, user.ID)

	if saving.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "saving not found under user"})
		return
	}

	saving.Name = body.Name
	saving.Description = body.Description
	saving.DepositAmount = body.DepositAmount
	saving.DepositFrequency = body.DepositFrequency
	saving.TotalAmount = body.TotalAmount


	result := initalisers.DB.Save(&saving)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, saving)
}

func DeleteSaving(ctx *gin.Context) {
	savingID := ctx.Param("id")

	user, err := getUserFromContext(ctx)
	checkUserError(ctx, err)

	result := initalisers.DB.Delete(&models.Saving{}, "id = ? AND user_id = ?", savingID, user.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "could not find saving to delete"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "saving deleted successfully"})
}

func lazyUpdateSavings(savings *[]models.Saving) {
	for i := range *savings {
		saving := &(*savings)[i]
		depositsOccured := calculateDepositsOccured(*saving)

		if depositsOccured > 0 {
			saving.TotalAmount = saving.TotalAmount + (depositsOccured * saving.DepositAmount)
			saving.LastCalculatedAt = time.Now()
			initalisers.DB.Save(&saving)
		}
	}
}

func calculateDepositsOccured(saving models.Saving) int {
	var depositsOccured int

	elapsedDays := math.Floor(time.Now().Sub(saving.LastCalculatedAt).Hours() / 24)
	switch saving.DepositFrequency {
	case "weekly":
		depositsOccured = int(math.Floor(elapsedDays / 7))
		break;
	case "fortnightly":
		depositsOccured = int(math.Floor(elapsedDays / 14))
		break;
	case "monthly":
		depositsOccured = int(math.Floor(elapsedDays / 28))
		break;
	}

	return depositsOccured
}
