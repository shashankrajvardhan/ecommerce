package main

import "time"

// Users
type Students struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
type Auth struct {
	Token string `json:"token"`
}
type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}
type Users struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Mobile    int    `json:"mobile"`
	AltMobile int    `json:"alt_mob"`
	Gender    string `json:"gender"`
}
type Address struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	Type            string `json:"type"`
	Address         string `json:"address"`
	Address2        string `json:"address2"`
	Address3        string `json:"address3"`
	ReceiversName   string `json:"receivers_name"`
	ReceiversMobile int    `json:"receivers_mobile"`
}

//Customer
type Customer struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

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
// "DELETE FROM cart WHERE user_id = $1 AND product_id = $2", a.UserID, a.ProductID