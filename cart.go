package main

import (
	"database/sql"
	cartfunc "ecommerce/packages"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var a Cart
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := cartfunc.Add(cartfunc.Cart{}, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully removed from cart"})
}

func RemoveItem(c *gin.Context) {
	var a struct {
		UserID    int    `json:"user_id"`
		ProductID int    `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Types     string `json:"types"`
	}

	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := cartfunc.Remove(cartfunc.Cart{}, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully removed from cart"})
}

func BuyFromCart(c *gin.Context) {
	var a struct {
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	var db *sql.DB
	if err := cartfunc.Buy(a.UserID, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully purchased from cart"})
}

func getProducts(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, price, quantity FROM product")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func Products(c *gin.Context) {
	products := []Product{
		{Name: "Electronics", Price: 100, Quantity: 10},
		{Name: "Cloths", Price: 200, Quantity: 20},
		{Name: "Grocery", Price: 300, Quantity: 30},
	}

	for _, p := range products {
		_, err := db.Exec("INSERT INTO product (name, price, quantity) VALUES ($1, $2, $3)", p.Name, p.Price, p.Quantity)
		if err != nil {
			log.Fatal("Failed to preload products:", err)
		}
	}
	c.JSON(http.StatusOK, products)
}
