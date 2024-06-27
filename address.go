package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func insert_address(c *gin.Context) {
	var a Address
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := db.QueryRow("INSERT INTO addresses (user_id, type, address, address2, address3, receivers_name, receivers_mobile) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", a.UserID, a.Type, a.Address, a.Address2, a.Address3, a.ReceiversName, a.ReceiversMobile).Scan(&a.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func update_address(c *gin.Context) {
	var a Address
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	aID := c.Param("id")
	err := db.QueryRow("UPDATE addresses SET user_id = $2, type = $3, address = $4, address2 = $5, address3 = $6, receivers_name = $7, receivers_mobile = $8 WHERE id = $1 RETURNING id", aID, a.UserID, a.Type, a.Address, a.Address2, a.Address3, a.ReceiversName, a.ReceiversMobile).Scan(&aID)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, a)
}
