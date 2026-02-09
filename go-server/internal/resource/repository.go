package resource

// Generic CRUD operations for resources (employees, projects, orders)

import (
	"fmt"
	"server/internal/config"
	"strings"
)

type Repository struct{}

// Helper to fetch allowed fields for a role/resource
func (r *Repository) getAllowedFields(roleID int, resource string) (viewFields map[string]bool, editFields map[string]bool, err error) {
	viewFields = make(map[string]bool)
	editFields = make(map[string]bool)

	// Admin (Role 1) has full access
	if roleID == 1 {
		return nil, nil, nil // nil maps signal full access
	}

	query := `
		SELECT rf.field_name, COALESCE(rfp.can_view, false), COALESCE(rfp.can_edit, false)
		FROM role_field_permissions rfp
		JOIN resource_fields rf ON rfp.resource_field_id = rf.id
		JOIN resources res ON rf.resource_id = res.id
		WHERE rfp.role_id = $1 AND res.name = $2
	`
	rows, err := config.DB.Query(query, roleID, resource)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var field string
		var view, edit bool
		if err := rows.Scan(&field, &view, &edit); err != nil {
			// Ensure we don't silently fail, but for now continue is safe if COALESCE works
			continue
		}
		if view {
			viewFields[field] = true
		}
		if edit {
			editFields[field] = true
		}
	}
	return viewFields, editFields, nil
}

func (r *Repository) GetAll(resource string, roleID int) ([]map[string]interface{}, error) {
	viewFields, _, err := r.getAllowedFields(roleID, resource)
	if err != nil {
		return nil, err
	}

	selectClause := "*"
	if roleID != 1 {
		cols := []string{"id"} // Always include ID
		// Retrieve all potential columns first to know what exists (or trust the allowed list)
		// For simplicity/safety, we only select fields that are explicitly allowed.
		// However, we need to know the valid columns of the table to avoid SQL errors if permission exists but column doesn't (rare).
		// Better approach: Select * then filter in Go? No, requirement says "Fields ... must NOT appear".
		// Construct SELECT based on viewFields.
		if len(viewFields) > 0 {
			for f := range viewFields {
				cols = append(cols, f)
			}
			selectClause = strings.Join(cols, ", ")
		} else {
			// If no fields allowed (but table access exists), return just IDs? or empty?
			// Requirement says "See only tables they have access to". If table access is yes but field access is none, show IDs?
			// Let's assume just ID if map is empty but not nil.
			selectClause = "id"
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", selectClause, resource)
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	results := []map[string]interface{}{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		row := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}
	return results, nil
}

func (r *Repository) Create(resource string, data map[string]interface{}, roleID int) (int, error) {
	_, editFields, err := r.getAllowedFields(roleID, resource)
	if err != nil {
		return 0, err
	}

	// Build dynamic INSERT
	keys := []string{}
	vals := []interface{}{}
	placeholders := []string{}
	i := 1

	for k, v := range data {
		// If strict permission check:
		if roleID != 1 {
			if !editFields[k] {
				continue // Skip fields user cannot edit (silently ignore)
			}
		}

		keys = append(keys, k)
		vals = append(vals, v)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	// Handle case where no fields are permitted (or data is empty)
	if len(keys) == 0 {
		// Try to insert with defaults (will fail if NOT NULL columns are missing, which is expected)
		query := fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING id", resource)
		var id int
		err = config.DB.QueryRow(query).Scan(&id)
		return id, err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		resource,
		strings.Join(keys, ", "),
		strings.Join(placeholders, ", "),
	)

	var id int
	err = config.DB.QueryRow(query, vals...).Scan(&id)
	return id, err
}

func (r *Repository) Update(resource, id string, data map[string]interface{}, roleID int) error {
	_, editFields, err := r.getAllowedFields(roleID, resource)
	if err != nil {
		return err
	}

	// Build dynamic UPDATE
	sets := []string{}
	vals := []interface{}{}
	i := 1

	for k, v := range data {
		// If strict permission check:
		if roleID != 1 {
			if !editFields[k] {
				continue // Skip fields user cannot edit (silently ignore)
			}
		}

		sets = append(sets, fmt.Sprintf("%s = $%d", k, i))
		vals = append(vals, v)
		i++
	}

	if len(sets) == 0 {
		return nil // Nothing to update
	}

	vals = append(vals, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		resource,
		strings.Join(sets, ", "),
		i,
	)

	_, err = config.DB.Exec(query, vals...)
	return err
}

func (r *Repository) Delete(resource, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", resource)
	_, err := config.DB.Exec(query, id)
	return err
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
