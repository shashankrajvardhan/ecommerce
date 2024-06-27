package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func insert_users(c *gin.Context) {
	var user Users
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := db.QueryRow("INSERT INTO users (name, mobile, alt_mob, gender) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Mobile, user.AltMobile, user.Gender).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func update_users(c *gin.Context) {
	var user Users
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	userID := c.Param("id")
	err := db.QueryRow("UPDATE users SET name = $2, mobile = $3, alt_mob = $4, gender = $5 WHERE id = $1 RETURNING id", userID, user.Name, user.Mobile, user.AltMobile, user.Gender).Scan(&user.ID)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
