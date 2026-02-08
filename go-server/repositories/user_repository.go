package repositories

import (
	"server/database"
	"server/models"

	"github.com/gin-gonic/gin"
)

type UserRepository struct{}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password, role_id, invitation_token, status FROM users WHERE username = $1"
	row := database.DB.QueryRow(query, username)

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.RoleID, &user.InvitationToken, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, password FROM users WHERE email = $1"
	// Note: We only need partial fields for checks in existing logic
	row := database.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user models.User) (int, error) {
	var id int
	query := "INSERT INTO users (username, password, role_id, status) VALUES ($1, $2, $3, 'Active') RETURNING id"
	err := database.DB.QueryRow(query, user.Username, user.Password, user.RoleID).Scan(&id)
	return id, err
}

func (r *UserRepository) CreateInvitedUser(email string, roleID int, token string) error {
	_, err := database.DB.Exec("INSERT INTO users (username, email, role_id, invitation_token, status) VALUES ($1, $2, $3, $4, 'Pending')", email, email, roleID, token)
	return err
}

func (r *UserRepository) UpdateInvitedUser(id int, roleID int, token string) error {
	_, err := database.DB.Exec("UPDATE users SET invitation_token = $1, role_id = $2, status = 'Pending' WHERE id = $3", token, roleID, id)
	return err
}

func (r *UserRepository) GetUserByToken(token string) (int, error) {
	var id int
	err := database.DB.QueryRow("SELECT id FROM users WHERE invitation_token = $1", token).Scan(&id)
	return id, err
}

func (r *UserRepository) ActivateUser(id int, password string) error {
	_, err := database.DB.Exec("UPDATE users SET password = $1, invitation_token = NULL, status = 'Active' WHERE id = $2", password, id)
	return err
}

func (r *UserRepository) DeclineUser(id int) error {
	_, err := database.DB.Exec("UPDATE users SET status = 'Declined', invitation_token = NULL WHERE id = $1", id)
	return err
}

func (r *UserRepository) GetAllUsers() ([]gin.H, error) {
	rows, err := database.DB.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []gin.H
	for rows.Next() {
		var id int
		var username string
		if err := rows.Scan(&id, &username); err == nil {
			userList = append(userList, gin.H{"id": id, "username": username})
		}
	}
	return userList, nil
}

func (r *UserRepository) GetAllUsersDetailed() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query("SELECT id, username, role_id, status FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.Status)
		users = append(users, map[string]interface{}{
			"id": u.ID, "username": u.Username, "role_id": u.RoleID, "status": u.Status,
		})
	}
	return users, nil
}

func (r *UserRepository) UpdateUserRole(userID int, roleID int) error {
	_, err := database.DB.Exec("UPDATE users SET role_id = $1 WHERE id = $2", roleID, userID)
	return err
}

func (r *UserRepository) DeleteUser(userID string) error {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}
