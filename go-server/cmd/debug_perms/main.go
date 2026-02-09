package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:123456789107@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT rrp.role_id, res.name, rrp.can_view FROM role_resource_permissions rrp JOIN resources res ON rrp.resource_id = res.id WHERE rrp.role_id = 4")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Permissions for Role 4:")
	for rows.Next() {
		var roleID int
		var resName string
		var canView bool
		rows.Scan(&roleID, &resName, &canView)
		fmt.Printf("Role: %d, Resource: %s, CanView: %v\n", roleID, resName, canView)
	}
}
