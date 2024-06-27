package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func secureCustomer(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing auth header"})
		return
	}
	tokenString = tokenString[len("Bearer "):]
	err := verifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}
	c.String(http.StatusOK, "Welcome to secure customer area")
}
