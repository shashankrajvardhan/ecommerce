package main

import (
	"database/sql"
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

	var productQuantity int
	err := db.QueryRow("SELECT quantity FROM product WHERE id = $1", a.ProductID).Scan(&productQuantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive product quantity"})
		return
	}

	if productQuantity < a.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough quantity available"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var currentQuentity int
	err = tx.QueryRow("SELECT quantity FROM cart WHERE user_id = $1 AND product_id = $2", a.UserID, a.ProductID).Scan(&currentQuentity)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO cart (user_id, product_id, types, quantity, price) VALUES ($1, $2, $3, $4, $5)", a.UserID, a.ProductID, a.Types, a.Quantity, a.Price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
				return
			}
		} else {
			_, err = tx.Exec("UPDATE cart SET quantity = quantity + $1 WHERE user_id = $2 AND product_id = $3", a.Quantity, a.UserID, a.ProductID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
				return
			}
		}

		_, err = tx.Exec("UPDATE product SET quantity = quantity - $1 WHERE id =$2", a.Quantity, a.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quality"})
			return
		}

		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successfully added to cart"})
	}
}

func RemoveItem(c *gin.Context) {
	var a Cart
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var productQuantity int
	err := db.QueryRow("SELECT quantity FROM product WHERE id $1", a.ProductID).Scan(&productQuantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive product quantity"})
		return
	}

	if productQuantity > a.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity is wrong"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var currentQuentity int
	err = tx.QueryRow("SELECT quantity FROM cart user_id = $1 AND product_id = $2", a.UserID, a.ProductID).Scan(&currentQuentity)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = 
		}
	}
}

func BuyFromCart(c *gin.Context) {

	var a struct {
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	rows, err := tx.Query("SELECT product_id, quantity FROM cart WHERE user_id = $1", a.UserID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart item"})
		return
	}
	rows.Close()

	for rows.Next() {
		var (
			productID, quantity int
		)
		if err := rows.Scan(&productID, &quantity); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan cart item"})
			return
		}

		_, err = tx.Exec("UPDATE product SET quantity = quantity - $1 WHERE id = $2", quantity, productID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product inventory"})
			return
		}
	}

	_, err = tx.Exec("UPDATE cart SET status = 'purchased' WHERE user_id = $1", a.UserID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
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
