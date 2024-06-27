package main

import "github.com/gin-gonic/gin"

func router() {

	r := gin.Default()
	// Users
	r.POST("/createuser", createUser)
	r.POST("/getbyusername", getByUsername)
	r.POST("/logout", logout)

	secureArea := r.Group("/", authMiddleware())

	secureArea.POST("/secure", secure)
	secureArea.POST("/insertusers", insert_users)
	secureArea.POST("/updateusers/:id", update_users)
	secureArea.POST("/insertaddress", insert_address)
	secureArea.POST("/updateaddress/:id", update_address)
	secureArea.POST("/addtocart", AddToCart)
	secureArea.POST("/removeitem", RemoveItem)
	secureArea.POST("/buyfromcart", BuyFromCart)
	secureArea.POST("/products", getProducts)
	secureArea.POST("/preload-products", Products)

	// Customer
	r.POST("/createcustomer", createCustomer)
	r.POST("/getByCustomer", getByCustomer)
	r.POST("/securecustomer", secureCustomer)
	r.Run(":8080")
}
