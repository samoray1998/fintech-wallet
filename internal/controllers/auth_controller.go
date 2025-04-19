package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samoray1998/fintech-wallet/internal/models"
	"github.com/samoray1998/fintech-wallet/internal/services"
)

type AuthController struct {
	authService *services.AuthService
	userService *services.UserServices
}

func NewAuthController(authService *services.AuthService, userService *services.UserServices) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var newUser struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		FullName: newUser.FullName,
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	createdUser, err := c.userService.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var creds struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.VerifyCredentials(creds.Email, creds.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials "})
		return
	}

	token, err := c.authService.GenerateTokens(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID.Hex(),
			"email":     user.Email,
			"kycStatus": user.KYCStatus,
		},
	})
}
