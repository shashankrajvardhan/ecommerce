package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func unhash(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func usernameExists(username string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM storage WHERE username = $1)", username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false
	}
	return exists
}
