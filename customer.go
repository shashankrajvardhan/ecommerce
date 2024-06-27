package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func createCustomer(c *gin.Context) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessionID := "u123"
	session[sessionID] = true
	c.SetCookie("session_id", sessionID, 60, "/", "", false, false)

	var s Customer
	if err := c.ShouldBind(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if usernameExists(s.Username) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username taken"})
		return
	}

	Password, err := hashPassword(s.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	token, err := createToken(s.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	_, err = db.Exec("INSERT INTO customer(username, password, token) VALUES ($1, $2, $3) RETURNING id", s.Username, Password, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	response := []gin.H{
		{"message": "User created successfully", "success": true},
		{"token": token},
	}
	c.JSON(http.StatusCreated, response)
}
