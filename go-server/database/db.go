package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost/rbac_db?sslmode=disable"
	}

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to Database: ", err)
	}

	fmt.Println("Connected to the Postgres database.")
	initializeSchema()
}

func initializeSchema() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE,
			password TEXT,
			role_id INTEGER REFERENCES roles(id),
			email TEXT,
			invitation_token TEXT,
			status TEXT DEFAULT 'Pending'
		)`,
		`CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			role_id INTEGER REFERENCES roles(id),
			resource TEXT,
			action TEXT,
			attributes TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			name TEXT,
			position TEXT,
			salary INTEGER,
			department TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name TEXT,
			assigned_to TEXT,
			status TEXT,
			budget INTEGER
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			customer_name TEXT,
			amount INTEGER,
			status TEXT,
			order_date TEXT
		)`,
	}

	for _, query := range tables {
		_, err := DB.Exec(query)
		if err != nil {
			log.Printf("Error creating table: %v\nQuery: %s", err, query)
		}
	}

	seedData()
}

func seedData() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM roles WHERE name = 'Admin'").Scan(&count)
	if err != nil {
		log.Println("Error checking admin role:", err)
		return
	}

	if count == 0 {
		var roleID int
		err = DB.QueryRow("INSERT INTO roles (name) VALUES ('Admin') RETURNING id").Scan(&roleID)
		if err != nil {
			log.Println("Error seeding Admin role:", err)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		_, err = DB.Exec("INSERT INTO users (username, password, role_id, status) VALUES ($1, $2, $3, 'Active')", "admin", string(hash), roleID)
		if err != nil {
			log.Println("Error seeding Admin user:", err)
		} else {
			fmt.Println("Admin user created")
		}

		resources := []string{"employees", "projects", "orders", "roles", "users"}
		actions := []string{"read", "create", "update", "delete"}

		for _, res := range resources {
			for _, action := range actions {
				_, err = DB.Exec("INSERT INTO permissions (role_id, resource, action, attributes) VALUES ($1, $2, $3, '*')", roleID, res, action)
				if err != nil {
					log.Println("Error seeding permissions:", err)
				}
			}
		}
		
		DB.Exec("INSERT INTO roles (name) VALUES ('user')")
	}
}
