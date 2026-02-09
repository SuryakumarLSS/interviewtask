package user

import "time"

type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Password        *string   `json:"password,omitempty"`
	RoleID          int       `json:"role_id"`
	Email           *string   `json:"email,omitempty"`
	InvitationToken *string   `json:"invitation_token,omitempty"`
	Status          string    `json:"status"`
	IsAdmin         bool      `json:"is_admin"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

type UpdateUserRoleRequest struct {
	RoleID int `json:"role_id"`
}
