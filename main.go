package main

import (
	"database/sql"
	conn "ecommerce/utils/dbconnection"
)

var db *sql.DB

func main() {
	conn.Connection()
	router()
}
