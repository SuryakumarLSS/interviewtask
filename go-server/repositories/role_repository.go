package repositories

import (
	"server/database"
	"server/models"
)

type RoleRepository struct{}

func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
	rows, err := database.DB.Query("SELECT id, name FROM roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err == nil {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

func (r *RoleRepository) CreateRole(name string) (int, error) {
	var id int
	err := database.DB.QueryRow("INSERT INTO roles (name) VALUES ($1) RETURNING id", name).Scan(&id)
	return id, err
}

func (r *RoleRepository) GetPermissions(roleID string) ([]models.Permission, error) {
	// Handling both int (from middleware) and string (from param)
	rows, err := database.DB.Query("SELECT id, role_id, resource, action, attributes FROM permissions WHERE role_id = $1", roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		rows.Scan(&p.ID, &p.RoleID, &p.Resource, &p.Action, &p.Attributes)
		perms = append(perms, p)
	}
	return perms, nil
}

// GetPermissionsByIntID avoids string conversion issues for internal calls
func (r *RoleRepository) GetPermissionsByIntID(roleID int) ([]models.Permission, error) {
	rows, err := database.DB.Query("SELECT id, role_id, resource, action, attributes FROM permissions WHERE role_id = $1", roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []models.Permission
	for rows.Next() {
		var p models.Permission
		rows.Scan(&p.ID, &p.RoleID, &p.Resource, &p.Action, &p.Attributes)
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *RoleRepository) CheckPermissionExists(roleID int, resource string, action string) (int, error) {
	var id int
	err := database.DB.QueryRow("SELECT id FROM permissions WHERE role_id = $1 AND resource = $2 AND action = $3", roleID, resource, action).Scan(&id)
	return id, err
}

func (r *RoleRepository) SetFullAccess(permID int) error {
	_, err := database.DB.Exec("UPDATE permissions SET attributes = '*' WHERE id = $1", permID)
	return err
}

func (r *RoleRepository) CreatePermission(roleID int, resource string, action string) (int, error) {
	var id int
	err := database.DB.QueryRow("INSERT INTO permissions (role_id, resource, action, attributes) VALUES ($1, $2, $3, '*') RETURNING id", roleID, resource, action).Scan(&id)
	return id, err
}

func (r *RoleRepository) DeletePermission(roleID string, resource string, action string) error {
	_, err := database.DB.Exec("DELETE FROM permissions WHERE role_id = $1 AND resource = $2 AND action = $3", roleID, resource, action)
	return err
}

func (r *RoleRepository) CheckAccess(roleID int, resource string, action string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM permissions 
			  WHERE role_id = $1 AND resource = $2 AND (action = $3 OR action = '*')`
	err := database.DB.QueryRow(query, roleID, resource, action).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
