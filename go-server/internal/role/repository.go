package role

import (
	"fmt"
	"server/internal/config"
)

type Repository struct{}

func (r *Repository) GetAll() ([]Role, error) {
	rows, err := config.DB.Query("SELECT id, name FROM roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		rows.Scan(&role.ID, &role.Name)
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *Repository) Create(name string) (int, error) {
	var id int
	err := config.DB.QueryRow("INSERT INTO roles (name) VALUES ($1) RETURNING id", name).Scan(&id)
	return id, err
}

func (r *Repository) Delete(id int) error {
	// Check if any users are using this role
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete role: it is assigned to %d users", count)
	}

	// Delete from legacy permissions table (doesn't have ON DELETE CASCADE)
	_, err = config.DB.Exec("DELETE FROM permissions WHERE role_id = $1", id)
	if err != nil {
		return err
	}

	// Delete the role (role_resource_permissions and role_field_permissions will cascade)
	_, err = config.DB.Exec("DELETE FROM roles WHERE id = $1", id)
	return err
}

func (r *Repository) GetPermissions(roleID string) ([]Permission, error) {
	// Query the new RBAC schema (role_resource_permissions table)
	query := `
		SELECT rrp.role_id, res.name, rrp.can_view, rrp.can_create, rrp.can_update, rrp.can_delete
		FROM role_resource_permissions rrp
		JOIN resources res ON rrp.resource_id = res.id
		WHERE rrp.role_id = $1
	`

	rows, err := config.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := []Permission{}
	for rows.Next() {
		var rID int
		var resourceName string
		var canView, canCreate, canUpdate, canDelete bool

		if err := rows.Scan(&rID, &resourceName, &canView, &canCreate, &canUpdate, &canDelete); err != nil {
			continue
		}

		// Convert to action-based format that frontend expects
		if canView {
			perms = append(perms, Permission{RoleID: rID, Resource: resourceName, Action: "read"})
		}
		if canCreate {
			perms = append(perms, Permission{RoleID: rID, Resource: resourceName, Action: "create"})
		}
		if canUpdate {
			perms = append(perms, Permission{RoleID: rID, Resource: resourceName, Action: "update"})
		}
		if canDelete {
			perms = append(perms, Permission{RoleID: rID, Resource: resourceName, Action: "delete"})
		}
	}
	return perms, nil
}

func (r *Repository) AddOrUpdatePermission(roleID int, resource, action string) (string, int, error) {
	// 1. Get resource ID
	var resourceID int
	err := config.DB.QueryRow("SELECT id FROM resources WHERE name = $1", resource).Scan(&resourceID)
	if err != nil {
		return "", 0, err
	}

	// 2. Map action to column
	column := ""
	switch action {
	case "read":
		column = "can_view"
	case "create":
		column = "can_create"
	case "update":
		column = "can_update"
	case "delete":
		column = "can_delete"
	}

	if column == "" {
		return "", 0, fmt.Errorf("invalid action")
	}

	// 3. Upsert into role_resource_permissions
	query := "INSERT INTO role_resource_permissions (role_id, resource_id, " + column + ") VALUES ($1, $2, TRUE) " +
		"ON CONFLICT (role_id, resource_id) DO UPDATE SET " + column + " = TRUE RETURNING id"

	var id int
	err = config.DB.QueryRow(query, roleID, resourceID).Scan(&id)
	return "created", id, err
}

func (r *Repository) DeletePermission(roleID, resource, action string) error {
	// 1. Get resource ID
	var resourceID int
	err := config.DB.QueryRow("SELECT id FROM resources WHERE name = $1", resource).Scan(&resourceID)
	if err != nil {
		return err
	}

	// 2. Map action to column
	column := ""
	switch action {
	case "read":
		column = "can_view"
	case "create":
		column = "can_create"
	case "update":
		column = "can_update"
	case "delete":
		column = "can_delete"
	}

	if column == "" {
		return nil
	}

	// 3. Update to FALSE
	_, err = config.DB.Exec("UPDATE role_resource_permissions SET "+column+" = FALSE WHERE role_id = $1 AND resource_id = $2", roleID, resourceID)
	return err
}

// GetFieldPermissions fetches specific field permissions for a role (including defaults)
func (r *Repository) GetFieldPermissions(roleID int) ([]FieldPermission, error) {
	query := `
		SELECT 
			r.name, 
			rf.field_name, 
			COALESCE(rfp.can_view, false), 
			COALESCE(rfp.can_edit, false)
		FROM resource_fields rf
		JOIN resources r ON rf.resource_id = r.id
		LEFT JOIN role_field_permissions rfp ON rf.id = rfp.resource_field_id AND rfp.role_id = $1
		ORDER BY r.name, rf.id;
	`
	rows, err := config.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []FieldPermission
	for rows.Next() {
		var p FieldPermission
		p.RoleID = roleID
		if err := rows.Scan(&p.Resource, &p.Field, &p.CanView, &p.CanEdit); err != nil {
			continue
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *Repository) UpdateFieldPermission(roleID int, resource, field string, canView, canEdit bool) error {
	// Nested subquery to find usage of resource_field_id
	query := `
		INSERT INTO role_field_permissions (role_id, resource_field_id, can_view, can_edit)
		VALUES (
			$1, 
			(SELECT rf.id FROM resource_fields rf JOIN resources r ON rf.resource_id = r.id WHERE r.name = $2 AND rf.field_name = $3), 
			$4, 
			$5
		)
		ON CONFLICT (role_id, resource_field_id) 
		DO UPDATE SET can_view = $4, can_edit = $5
	`
	_, err := config.DB.Exec(query, roleID, resource, field, canView, canEdit)
	return err
}
