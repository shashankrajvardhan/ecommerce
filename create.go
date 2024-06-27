package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	session      = make(map[string]bool)
	sessionMutex = &sync.Mutex{}
)

func createUser(c *gin.Context) {

	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	sessionID := "u123"
	session[sessionID] = true
	c.SetCookie("session_id", sessionID, 60, "/", "", false, false)

	var u Students
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if usernameExists(u.Username) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username taken"})
		return
	}

	Password, err := hashPassword(u.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Token
	token, err := createToken(u.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	_, err = db.Exec("INSERT INTO storage(username, password, token) VALUES ($1, $2, $3) RETURNING id", u.Username, Password, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	responses := []gin.H{
		{"message": "User created successfully", "success": true},
		{"token": token},
	}
	c.JSON(http.StatusCreated, responses)
}
