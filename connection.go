package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func connection() *sql.DB {
	dsn := "host=localhost user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected")
	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS storage (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			token VARCHAR (2550) NOT NULL
			);
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Created Successfully")

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			mobile BIGINT NOT NULL,
			alt_mob BIGINT NOT NULL,
			gender VARCHAR(255) NOT NULL
	);
		`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Users Created Successfully")

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS addresses (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			type VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			address2 VARCHAR(255) NOT NULL,
			address3 VARCHAR(255) NOT NULL,
			receivers_name VARCHAR NOT NULL,
			receivers_mobile BIGINT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
	);
		`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Addresses Created Successfully")

	// Customer Table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS customer (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		token VARCHAR (2550) NOT NULL
		);
`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Customer Created Successfully")

	// Cart
	_, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS cart (
	  id SERIAL PRIMARY KEY,
	  user_id INT NOT NULL,
	  product_id INT NOT NULL,
	  types VARCHAR(50) NOT NULL,
	  quantity INT NOT NULL,
	  price INT NOT NULL,
	  status VARCHAR(50) DEFAULT 'pending',
	  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	 
    );
`)
	if err != nil {
		log.Fatal("Failed to create table", err)
	}
	fmt.Println("Cart created Successfully")

	// Product
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS product (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	price INT NOT NULL,
	quantity INT NOT NULL
	);
	`)
	if err != nil{
		log.Fatal("Failed to create table", err)
	}
	fmt.Println("Product created Successfully")
	return db
}
