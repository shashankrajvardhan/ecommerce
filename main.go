package main

import "database/sql"

var db *sql.DB

func main() {
	connection()
	router()
}
