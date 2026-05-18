package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/vyonyx/at-what-cost/initalisers"
	"github.com/vyonyx/at-what-cost/models"
)

type SignUpDetails struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginDetails struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func Signup(ctx *gin.Context) {
	var body SignUpDetails
	err := ctx.ShouldBind(&body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not create user"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not hash password",
		})
		return
	}

	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash)}

	result := initalisers.DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func Login(ctx *gin.Context)  {
	var body LoginDetails
	if ctx.ShouldBind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read credentials"})
	}

	var user models.User
	initalisers.DB.Find(&user, "email = ?", body.Email)

	if  user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or password"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create JWT token",
		})

		return
	}

	// Respond
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"Validation": tokenString})
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("Authorization", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Loggout out successfully",
	})
}

func Validate( ctx *gin.Context)  {
	ctx.JSON(http.StatusOK, gin.H{"message": "user login validated"})
}
