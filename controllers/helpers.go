package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vyonyx/at-what-cost/models"
)

var (
	userNotSignedInError = errors.New("user not signed in")
)

func getUserFromContext(ctx *gin.Context) (models.User, error) {
	details, userExists := ctx.Get("user")

	if !userExists {
		return models.User{}, userNotSignedInError
	}

	user := details.(models.User)
	return user, nil
}

func checkUserError(ctx *gin.Context, err error) {
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
}
