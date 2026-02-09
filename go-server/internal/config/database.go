package config

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
		// Core Auth Tables
		`CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// 		INSERT INTO roles (name) VALUES
		// ('viewer'),
		// ('editor'),
		// ('manager');

		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT,
			role_id INTEGER REFERENCES roles(id),
			email TEXT,
			invitation_token TEXT,
			status TEXT DEFAULT 'Pending',
			is_admin BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// 		-- Admin
		// INSERT INTO users (username, password_hash, is_admin)
		// VALUES ('admin', 'admin123', TRUE);

		// -- Regular users
		// INSERT INTO users (username, password_hash, role_id)
		// VALUES
		// ('viewer1', 'viewer123', 1),
		// ('editor1', 'editor123', 2),
		// ('manager1', 'manager123', 3);

		// Resource Metadata (Config-Driven Design)
		`CREATE TABLE IF NOT EXISTS resources (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			display_name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// 	INSERT INTO resources (name) VALUES
		// ('employees'),
		// ('projects'),
		// ('orders');

		`CREATE TABLE IF NOT EXISTS resource_fields (
			id SERIAL PRIMARY KEY,
			resource_id INTEGER REFERENCES resources(id) ON DELETE CASCADE,
			field_name TEXT NOT NULL,
			data_type TEXT DEFAULT 'text',
			is_sensitive BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(resource_id, field_name)
		)`,

		// Permission Tables (Table-Level & Field-Level)
		`CREATE TABLE IF NOT EXISTS role_resource_permissions (
			id SERIAL PRIMARY KEY,
			role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
			resource_id INTEGER REFERENCES resources(id) ON DELETE CASCADE,
			can_view BOOLEAN DEFAULT FALSE,
			can_create BOOLEAN DEFAULT FALSE,
			can_update BOOLEAN DEFAULT FALSE,
			can_delete BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(role_id, resource_id)
		)`,

		`CREATE TABLE IF NOT EXISTS role_field_permissions (
			id SERIAL PRIMARY KEY,
			role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
			resource_field_id INTEGER REFERENCES resource_fields(id) ON DELETE CASCADE,
			can_view BOOLEAN DEFAULT FALSE,
			can_edit BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(role_id, resource_field_id)
		)`,

		// Legacy permissions table (kept for backward compatibility, will be deprecated)
		`CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			role_id INTEGER REFERENCES roles(id),
			resource TEXT,
			action TEXT,
			attributes TEXT
		)`,

		// Business Data Tables
		`CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			position TEXT,
			salary INTEGER,
			department TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			assigned_to TEXT,
			status TEXT,
			budget INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			customer_name TEXT NOT NULL,
			amount INTEGER,
			status TEXT,
			order_date TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	// Check if Admin role exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM roles WHERE name = 'Admin'").Scan(&count)
	if err != nil {
		log.Println("Error checking admin role:", err)
		return
	}

	if count == 0 {
		// Create Admin role
		var roleID int
		err = DB.QueryRow("INSERT INTO roles (name) VALUES ('Admin') RETURNING id").Scan(&roleID)
		if err != nil {
			log.Println("Error seeding Admin role:", err)
			return
		}

		// Create Admin user with is_admin flag
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)
		_, err = DB.Exec(
			"INSERT INTO users (username, password, role_id, status, is_admin) VALUES ($1, $2, $3, 'Active', TRUE)",
			"admin", string(hash), roleID,
		)
		if err != nil {
			log.Println("Error seeding Admin user:", err)
		} else {
			fmt.Println("✅ Superadmin user created (username: admin, password: admin123)")
		}

		// Create default "User" role
		DB.Exec("INSERT INTO roles (name) VALUES ('User')")

		// Seed Resources
		resourceNames := []string{"employees", "projects", "orders", "roles", "users"}
		resourceDisplayNames := map[string]string{
			"employees": "Employees",
			"projects":  "Projects",
			"orders":    "Orders",
			"roles":     "Roles",
			"users":     "Users",
		}

		for _, resName := range resourceNames {
			var resID int
			err = DB.QueryRow(
				"INSERT INTO resources (name, display_name) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET display_name = $2 RETURNING id",
				resName, resourceDisplayNames[resName],
			).Scan(&resID)
			if err != nil {
				log.Printf("Error seeding resource %s: %v", resName, err)
				continue
			}

			// Seed fields for each resource
			seedFieldsForResource(resID, resName)

			// Grant full table-level permissions to Admin role
			_, err = DB.Exec(
				`INSERT INTO role_resource_permissions (role_id, resource_id, can_view, can_create, can_update, can_delete)
				 VALUES ($1, $2, TRUE, TRUE, TRUE, TRUE)
				 ON CONFLICT (role_id, resource_id) DO UPDATE SET can_view = TRUE, can_create = TRUE, can_update = TRUE, can_delete = TRUE`,
				roleID, resID,
			)
			if err != nil {
				log.Printf("Error seeding resource permissions: %v", err)
			}

			// Grant full field-level permissions to Admin role
			grantFieldPermissionsToAdmin(roleID, resID)
		}

		// Migrate legacy permissions (if any exist)
		migrateLegacyPermissions()

		fmt.Println("✅ Database seeded with resources, fields, and permissions")
	}
}

func seedFieldsForResource(resourceID int, resourceName string) {
	fieldMap := map[string][]map[string]interface{}{
		"employees": {
			{"name": "name", "type": "text", "sensitive": false},
			{"name": "position", "type": "text", "sensitive": false},
			{"name": "salary", "type": "number", "sensitive": true},
			{"name": "department", "type": "text", "sensitive": false},
		},
		"projects": {
			{"name": "name", "type": "text", "sensitive": false},
			{"name": "assigned_to", "type": "text", "sensitive": false},
			{"name": "status", "type": "text", "sensitive": false},
			{"name": "budget", "type": "number", "sensitive": true},
		},
		"orders": {
			{"name": "customer_name", "type": "text", "sensitive": false},
			{"name": "amount", "type": "number", "sensitive": false},
			{"name": "status", "type": "text", "sensitive": false},
			{"name": "order_date", "type": "text", "sensitive": false},
		},
		"roles": {
			{"name": "name", "type": "text", "sensitive": false},
		},
		"users": {
			{"name": "username", "type": "text", "sensitive": false},
			{"name": "email", "type": "text", "sensitive": true},
			{"name": "status", "type": "text", "sensitive": false},
		},
	}

	fields, ok := fieldMap[resourceName]
	if !ok {
		return
	}

	for _, field := range fields {
		_, err := DB.Exec(
			`INSERT INTO resource_fields (resource_id, field_name, data_type, is_sensitive)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (resource_id, field_name) DO NOTHING`,
			resourceID, field["name"], field["type"], field["sensitive"],
		)
		if err != nil {
			log.Printf("Error seeding field %s for %s: %v", field["name"], resourceName, err)
		}
	}
}

func grantFieldPermissionsToAdmin(roleID, resourceID int) {
	rows, err := DB.Query("SELECT id FROM resource_fields WHERE resource_id = $1", resourceID)
	if err != nil {
		log.Printf("Error fetching fields for resource %d: %v", resourceID, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var fieldID int
		rows.Scan(&fieldID)
		_, err = DB.Exec(
			`INSERT INTO role_field_permissions (role_id, resource_field_id, can_view, can_edit)
			 VALUES ($1, $2, TRUE, TRUE)
			 ON CONFLICT (role_id, resource_field_id) DO UPDATE SET can_view = TRUE, can_edit = TRUE`,
			roleID, fieldID,
		)
		if err != nil {
			log.Printf("Error seeding field permission: %v", err)
		}
	}
}

func migrateLegacyPermissions() {
	// This function can be used to migrate old permissions table to new structure if needed
	// For now, we'll keep both systems running in parallel
	log.Println("Legacy permissions table maintained for backward compatibility")
}
