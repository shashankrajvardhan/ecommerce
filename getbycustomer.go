package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getByCustomer(c *gin.Context) {
	var e Customer
	if err := c.ShouldBind(&e); err != nil {
		return
	}
	var f Customer

	err := db.QueryRow("SELECT * FROM customer WHERE username = $1", e.Username).Scan(&f.ID, &f.Username, &f.Password, &f.Token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	er := unhash([]byte(f.Password), []byte(e.Password))
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unhash password"})
		return
	}

	response := []gin.H{
		{"Successful": "Welcome", "success": true},
		{"token": f.Token},
	}
	c.JSON(http.StatusOK, response)
}
