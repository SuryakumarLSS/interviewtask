package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Hardcoded connection string from .env for debugging
	connStr := "postgres://postgres:123456789107@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id int
	var username string
	var passwordHash string
	var isAdmin bool

	// Check Admin user
	err = db.QueryRow("SELECT id, username, password, is_admin FROM users WHERE username = 'admin'").Scan(&id, &username, &passwordHash, &isAdmin)
	if err != nil {
		log.Fatalf("Error querying admin: %v\nRUNNING DB RESET...", err)
	}

	fmt.Printf("found Admin: ID=%d Username=%s Hash=%s IsAdmin=%v\n", id, username, passwordHash, isAdmin)

	// Verify Password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("admin123"))
	if err != nil {
		fmt.Printf("❌ Password verification FAILED: %v\n", err)

		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		_, err = db.Exec("UPDATE users SET password = $1 WHERE username = 'admin'", string(hash))
		if err != nil {
			fmt.Printf("Failed to reset password: %v\n", err)
		} else {
			fmt.Println("✅ Admin password manually reset to 'admin123'")
		}
	} else {
		fmt.Println("✅ Password verification SUCCESS (admin123 matches hash)")
	}
}
