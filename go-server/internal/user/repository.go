package user

import (
	"server/internal/config"
)

type Repository struct{}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password, role_id, invitation_token, status, is_admin FROM users WHERE username = $1"
	row := config.DB.QueryRow(query, username)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.RoleID, &user.InvitationToken, &user.Status, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByToken(token string) (*User, error) {
	var user User
	query := "SELECT id, username, role_id, invitation_token, status FROM users WHERE invitation_token = $1"
	row := config.DB.QueryRow(query, token)

	err := row.Scan(&user.ID, &user.Username, &user.RoleID, &user.InvitationToken, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(username, hashedPassword string, roleID int, email *string, invitationToken *string) (int, error) {
	var id int
	err := config.DB.QueryRow(
		"INSERT INTO users (username, password, role_id, email, invitation_token, status) VALUES ($1, $2, $3, $4, $5, 'Pending') RETURNING id",
		username, hashedPassword, roleID, email, invitationToken,
	).Scan(&id)
	return id, err
}

func (r *Repository) UpdateRole(userID, roleID int) error {
	_, err := config.DB.Exec("UPDATE users SET role_id = $1 WHERE id = $2", roleID, userID)
	return err
}

func (r *Repository) UpdatePassword(userID int, hashedPassword string) error {
	_, err := config.DB.Exec(
		"UPDATE users SET password = $1, invitation_token = NULL, status = 'Active' WHERE id = $2",
		hashedPassword, userID,
	)
	return err
}

func (r *Repository) Delete(userID string) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}

func (r *Repository) GetAll() ([]map[string]interface{}, error) {
	rows, err := config.DB.Query("SELECT id, username, role_id, status, is_admin FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.Status, &u.IsAdmin)
		users = append(users, map[string]interface{}{
			"id": u.ID, "username": u.Username, "role_id": u.RoleID, "status": u.Status, "is_admin": u.IsAdmin,
		})
	}
	return users, nil
}

func (r *Repository) IsAdmin(userID int) (bool, error) {
	var isAdmin bool
	err := config.DB.QueryRow("SELECT is_admin FROM users WHERE id = $1", userID).Scan(&isAdmin)
	return isAdmin, err
}

func (r *Repository) UsernameExists(username string) (bool, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	return count > 0, err
}
