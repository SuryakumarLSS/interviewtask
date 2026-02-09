package role

import "time"

type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type Permission struct {
	ID         int    `json:"id"`
	RoleID     int    `json:"role_id"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Attributes string `json:"attributes"`
}

type PermissionRequest struct {
	RoleID   int    `json:"role_id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type FieldPermission struct {
	RoleID   int    `json:"role_id"`
	Resource string `json:"resource"`
	Field    string `json:"field"`
	CanView  bool   `json:"can_view"`
	CanEdit  bool   `json:"can_edit"`
}

type UpdateFieldPermissionRequest struct {
	RoleID   int    `json:"role_id"`
	Resource string `json:"resource"`
	Field    string `json:"field"`
	CanView  bool   `json:"can_view"`
	CanEdit  bool   `json:"can_edit"`
}
