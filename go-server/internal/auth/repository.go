package auth

import (
	"server/internal/config"
)

type Repository struct{}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password, role_id, invitation_token, status, is_admin FROM users WHERE username = $1"
	row := config.DB.QueryRow(query, username)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.RoleID, &user.InvitationToken, &user.Status, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByToken(token string) (*User, error) {
	var user User
	query := "SELECT id, username, role_id, invitation_token, status FROM users WHERE invitation_token = $1"
	row := config.DB.QueryRow(query, token)

	err := row.Scan(&user.ID, &user.Username, &user.RoleID, &user.InvitationToken, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) SetPassword(userID int, hashedPassword string) error {
	_, err := config.DB.Exec(
		"UPDATE users SET password = $1, invitation_token = NULL, status = 'Active' WHERE id = $2",
		hashedPassword, userID,
	)
	return err
}

func (r *Repository) DeclineInvitation(userID int) error {
	_, err := config.DB.Exec("UPDATE users SET status = 'Declined', invitation_token = NULL WHERE id = $1", userID)
	return err
}

func (r *Repository) GetPermissionsByRoleID(roleID int) ([]Permission, error) {
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
		var roleID int
		var resourceName string
		var canView, canCreate, canUpdate, canDelete bool

		if err := rows.Scan(&roleID, &resourceName, &canView, &canCreate, &canUpdate, &canDelete); err != nil {
			continue
		}

		// Convert to action-based format that frontend expects
		if canView {
			perms = append(perms, Permission{
				RoleID:   roleID,
				Resource: resourceName,
				Action:   "read",
			})
		}
		if canCreate {
			perms = append(perms, Permission{
				RoleID:   roleID,
				Resource: resourceName,
				Action:   "create",
			})
		}
		if canUpdate {
			perms = append(perms, Permission{
				RoleID:   roleID,
				Resource: resourceName,
				Action:   "update",
			})
		}
		if canDelete {
			perms = append(perms, Permission{
				RoleID:   roleID,
				Resource: resourceName,
				Action:   "delete",
			})
		}
	}

	return perms, nil
}

type Permission struct {
	ID         int    `json:"id"`
	RoleID     int    `json:"role_id"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Attributes string `json:"attributes"`
}

func (r *Repository) GetAllUsersSimple() ([]map[string]interface{}, error) {
	rows, err := config.DB.Query("SELECT id, username FROM users WHERE status = 'Active'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username string
		rows.Scan(&id, &username)
		users = append(users, map[string]interface{}{"id": id, "username": username})
	}
	return users, nil
}

func (r *Repository) GetEmployeeNames() ([]map[string]interface{}, error) {
	rows, err := config.DB.Query("SELECT id, name FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		employees = append(employees, map[string]interface{}{"id": id, "name": name})
	}
	return employees, nil
}

// Helper to check role permission
func (r *Repository) GetFieldPermissionsByRoleID(roleID int) ([]FieldPermission, error) {
	var query string
	var args []interface{}

	if roleID == 1 {
		// Admin gets everything
		query = `
			SELECT res.name, rf.field_name, true, true
			FROM resource_fields rf
			JOIN resources res ON rf.resource_id = res.id
		`
	} else {
		query = `
			SELECT res.name, rf.field_name, rfp.can_view, rfp.can_edit
			FROM role_field_permissions rfp
			JOIN resource_fields rf ON rfp.resource_field_id = rf.id
			JOIN resources res ON rf.resource_id = res.id
			WHERE rfp.role_id = $1
		`
		args = append(args, roleID)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := []FieldPermission{}
	for rows.Next() {
		var p FieldPermission
		if err := rows.Scan(&p.Resource, &p.Field, &p.CanView, &p.CanEdit); err != nil {
			continue
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *Repository) HasPermission(roleID int, resource, action string) (bool, error) {
	query := `
		SELECT CASE 
			WHEN $3 = 'read' THEN rrp.can_view
			WHEN $3 = 'create' THEN rrp.can_create
			WHEN $3 = 'update' THEN rrp.can_update
			WHEN $3 = 'delete' THEN rrp.can_delete
			ELSE false
		END as has_permission
		FROM role_resource_permissions rrp
		JOIN resources res ON rrp.resource_id = res.id
		WHERE rrp.role_id = $1 AND res.name = $2
	`

	var hasPermission bool
	err := config.DB.QueryRow(query, roleID, resource, action).Scan(&hasPermission)
	if err != nil {
		return false, nil // No permission found
	}

	return hasPermission, nil
}
