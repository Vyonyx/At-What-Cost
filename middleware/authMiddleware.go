package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vyonyx/at-what-cost/initalisers"
	"github.com/vyonyx/at-what-cost/models"
)

func RequireAuth(ctx *gin.Context) {
	// Get cookie off request
	tokenString, err := ctx.Cookie("Authorization")

	if tokenString == "" || err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization token not attached on request to protected route",
		})

		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode/validate cookie
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check expiry date on cookie
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user with token sub
		var user models.User
		initalisers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach user details to context
		ctx.Set("user", user)

		// Continue
		ctx.Next()

	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
