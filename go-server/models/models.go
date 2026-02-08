package models

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	Password        *string `json:"password,omitempty"` // Pointer to handle NULL
	RoleID          int     `json:"role_id"`
	Email           *string `json:"email,omitempty"`
	InvitationToken *string `json:"invitation_token,omitempty"`
	Status          string  `json:"status"`
}

type Permission struct {
	ID         int    `json:"id"`
	RoleID     int    `json:"role_id"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Attributes string `json:"attributes"`
}

// Resource structs
type Employee struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	Salary     int    `json:"salary"`
	Department string `json:"department"`
}

type Project struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AssignedTo string `json:"assigned_to"`
	Status     string `json:"status"`
	Budget     int    `json:"budget"`
}

type Order struct {
	ID           int    `json:"id"`
	CustomerName string `json:"customer_name"`
	Amount       int    `json:"amount"`
	Status       string `json:"status"`
	OrderDate    string `json:"order_date"`
}
