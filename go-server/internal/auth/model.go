package auth

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type SetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type DeclineInvitationRequest struct {
	Token string `json:"token"`
}

// User struct for auth responses
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

type FieldPermission struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	CanView  bool   `json:"can_view"`
	CanEdit  bool   `json:"can_edit"`
}
