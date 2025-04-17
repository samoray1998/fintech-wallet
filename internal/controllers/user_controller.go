package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samoray1998/fintech-wallet/internal/services"
)

type UserController struct {
	userService *services.UserServices
}

func NewUserController(userService *services.UserServices) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := c.userService.GetUserByID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":        user.ID.Hex(),
		"email":     user.Email,
		"kycStatus": user.KYCStatus,
		"createdAt": user.CreatedAt,
	})
}

func (c *UserController) InitiateKYC(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// In a real app, you would process KYC documents here
	_, err := c.userService.UpdateKYCStatus(userID.(string), "pending")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "KYC verification initiated"})
}
