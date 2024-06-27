package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getByUsername(c *gin.Context) {
	var u Students
	if err := c.ShouldBind(&u); err != nil {
		return
	}
	var a Students

	err := db.QueryRow("SELECT * FROM storage WHERE username = $1", u.Username).Scan(&a.ID, &a.Username, &a.Password, &a.Token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	er := unhash([]byte(a.Password), []byte(u.Password))
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unhash password"})
		return
	}

	responses := []gin.H{
		{"Successful": "Welcome", "success": true},
		{"token": a.Token},
	}
	c.JSON(http.StatusCreated, responses)
}
