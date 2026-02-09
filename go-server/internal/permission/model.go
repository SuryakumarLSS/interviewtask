package permission

import "time"

// Resource metadata
type Resource struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type ResourceField struct {
	ID          int       `json:"id"`
	ResourceID  int       `json:"resource_id"`
	FieldName   string    `json:"field_name"`
	DataType    string    `json:"data_type"`
	IsSensitive bool      `json:"is_sensitive"`
	CreatedAt   time.Time `json:"created_at"`
}

// Table-level permissions
type RoleResourcePermission struct {
	ID         int       `json:"id"`
	RoleID     int       `json:"role_id"`
	ResourceID int       `json:"resource_id"`
	CanView    bool      `json:"can_view"`
	CanCreate  bool      `json:"can_create"`
	CanUpdate  bool      `json:"can_update"`
	CanDelete  bool      `json:"can_delete"`
	CreatedAt  time.Time `json:"created_at"`
}

// Field-level permissions
type RoleFieldPermission struct {
	ID              int       `json:"id"`
	RoleID          int       `json:"role_id"`
	ResourceFieldID int       `json:"resource_field_id"`
	CanView         bool      `json:"can_view"`
	CanEdit         bool      `json:"can_edit"`
	CreatedAt       time.Time `json:"created_at"`
}

// Legacy permission (backward compatibility)
type Permission struct {
	ID         int    `json:"id"`
	RoleID     int    `json:"role_id"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Attributes string `json:"attributes"`
}

// Request/Response DTOs
type SavePermissionRequest struct {
	RoleID   int    `json:"role_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type DeletePermissionRequest struct {
	RoleID   string `json:"role_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// Permission evaluation result
type AccessResult struct {
	Allowed        bool     `json:"allowed"`
	ViewableFields []string `json:"viewable_fields,omitempty"`
	EditableFields []string `json:"editable_fields,omitempty"`
}
