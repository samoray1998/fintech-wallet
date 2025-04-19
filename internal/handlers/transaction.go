package handlers

import "github.com/gin-gonic/gin"

func Transfer(c *gin.Context) {
	// 1. Check sender balance
	// 2. Deduct from sender, add to receiver
	// 3. Record transaction
	// 4. Return success/failure
}
