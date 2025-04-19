package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samoray1998/fintech-wallet/internal/controllers"
	"github.com/samoray1998/fintech-wallet/internal/middlewares"
)

func SetupRouter(
	authMiddleware *middlewares.AuthMiddleware,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	//accountController *controllers.AccountController,
	//transactionController *controllers.TransactionController,
	//rateController *controllers.RateController,
	rateLimit int,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	//router.Use(middlewares.RateLimiter(rateLimit))

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
		//public.GET("/rates", rateController.GetCurrentRates)
	}

	// Authenticated routes
	private := router.Group("/api/v1")
	private.Use(authMiddleware.Authenticate)
	{
		private.GET("/users/me", userController.GetProfile)
		//private.POST("/accounts", accountController.CreateAccount)
		//private.POST("/transactions", transactionController.CreateTransaction)
	}

	return router
}
