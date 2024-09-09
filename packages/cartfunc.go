package cartfunc

import (
	"database/sql"
	"errors"
	"time"
)

type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Types     string    `json:"types"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Add(a Cart, db *sql.DB) error {

	var productQuantity int
	err := db.QueryRow("SELECT quantity FROM product WHERE id = $1", a.ProductID).Scan(&productQuantity)
	if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive product quantity"})
		return errors.New("failed to retrive product quantity")
	}

	if productQuantity < a.Quantity {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough quantity available"})
		return errors.New("not enough quantity available")
	}

	tx, err := db.Begin()
	if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return errors.New("failed to start transaction")
	}

	var currentQuentity int
	err = tx.QueryRow("SELECT quantity FROM cart WHERE user_id = $1 AND product_id = $2", a.UserID, a.ProductID).Scan(&currentQuentity)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("INSERT INTO cart (user_id, product_id, types, quantity, price) VALUES ($1, $2, $3, $4, $5)", a.UserID, a.ProductID, a.Types, a.Quantity, a.Price)
			if err != nil {
				//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
				return errors.New("failed to add to cart")
			}
		} else {
			_, err = tx.Exec("UPDATE cart SET quantity = quantity + $1 WHERE user_id = $2 AND product_id = $3", a.Quantity, a.UserID, a.ProductID)
			if err != nil {
				//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
				return errors.New("failed to update cart")
			}
		}

		_, err = tx.Exec("UPDATE product SET quantity = quantity - $1 WHERE id =$2", a.Quantity, a.ProductID)
		if err != nil {
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quality"})
			return errors.New("failed to update product quality")
		}

		err = tx.Commit()
		if err != nil {
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return errors.New("failed to commit transaction")
		}
		//	c.JSON(http.StatusOK, gin.H{"message": "successfully added to cart"})
	}
	return nil
}

func Remove(a Cart, db *sql.DB) error {

	var productQuantity int
	err := db.QueryRow("SELECT quantity FROM product WHERE id $1", a.ProductID).Scan(&productQuantity)
	if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive product quantity"})
		return errors.New("failed to retrive product quantity")
	}

	if productQuantity > a.Quantity {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity is wrong"})
		return errors.New("quantity is wrong")
	}

	tx, err := db.Begin()
	if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return errors.New("failed to start transaction")
	}

	var currentQuentity int
	err = tx.QueryRow("SELECT quantity FROM cart user_id = $1 AND product_id = $2", a.UserID, a.ProductID).Scan(&currentQuentity)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = tx.Exec("DELETE FROM cart WHERE user_id = $1 AND product_id = $2 AND types = $3", a.UserID, a.ProductID, a.Types)
			if err != nil {
				//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from cart"})
				return errors.New("failed to remove from cart")
			}
		} else {
			_, err = tx.Exec("UPDATE cart SET quantity = quantity - $1 WHERE user_id = $2 AND product_id = $3", a.Quantity, a.UserID, a.ProductID)
			if err != nil {
				//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
				return errors.New("failed to update cart")
			}
		}

		_, err = tx.Exec("UPDATE product SET quantity = quantity + $1 WHERE id = $2", a.Quantity, a.ProductID)
		if err != nil {
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quality"})
			return errors.New("failed to update product quality")
		}

		err = tx.Commit()
		if err != nil {
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to do update"})
			return errors.New("failed to do update")
		}
		//	c.JSON(http.StatusOK, gin.H{"message": "successfully removed from cart"})
	}
	return nil
}

func Buy(userID int, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("failed to start transaction")
	}

	rows, err := tx.Query("SELECT product_id, quantity FROM cart WHERE user_id = $1", userID)
	if err != nil {
		tx.Rollback()
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart item"})
		return errors.New("failed to get cart item")
	}
	rows.Close()

	for rows.Next() {
		var (
			productID, quantity int
		)
		if err := rows.Scan(&productID, &quantity); err != nil {
			tx.Rollback()
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan cart item"})
			return errors.New("failed to scan cart item")
		}

		_, err = tx.Exec("UPDATE product SET quantity = quantity - $1 WHERE id = $2", quantity, productID)
		if err != nil {
			tx.Rollback()
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product inventory"})
			return errors.New("failed to update product inventory")
		}
	}
	if _, err = tx.Exec("UPDATE cart SET status = 'purchased' WHERE user_id = $1", userID); err != nil {
		tx.Rollback()
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return errors.New("failed to update cart")
	}
	if err := tx.Commit(); err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return errors.New("failed to commit transaction")
	}

	//	c.JSON(http.StatusOK, gin.H{"message": "Successfully purchased from cart"})
	return nil
}
