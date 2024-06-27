package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func logout(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully", "success": true})
}
