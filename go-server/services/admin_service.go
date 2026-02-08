package services

import (
	"database/sql"
	"fmt"
	"server/models"
	"server/repositories"
	"server/utils"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	UserRepo *repositories.UserRepository
	RoleRepo *repositories.RoleRepository
}

func (s *AdminService) GetRoles() ([]models.Role, error) {
	return s.RoleRepo.GetAllRoles()
}

func (s *AdminService) CreateRole(name string) (int, error) {
	return s.RoleRepo.CreateRole(name)
}

func (s *AdminService) GetPermissions(roleID string) ([]models.Permission, error) {
	return s.RoleRepo.GetPermissions(roleID)
}

func (s *AdminService) AddOrUpdatePermission(roleID int, resource, action string) (string, int, error) {
	existingID, err := s.RoleRepo.CheckPermissionExists(roleID, resource, action)
	if err == nil {
		// Exists
		err := s.RoleRepo.SetFullAccess(existingID)
		if err != nil {
			return "", 0, err
		}
		return "updated", 0, nil
	}

	// Create
	newID, err := s.RoleRepo.CreatePermission(roleID, resource, action)
	return "created", newID, err
}

func (s *AdminService) DeletePermission(roleID, resource, action string) error {
	return s.RoleRepo.DeletePermission(roleID, resource, action)
}

func (s *AdminService) GetAllUsers() ([]map[string]interface{}, error) {
	return s.UserRepo.GetAllUsersDetailed()
}

func (s *AdminService) CreateUserOrInvite(req struct {
	Username string
	Password string
	RoleID   int
	Email    string
}) (string, string, int, error) { // returns type(invite/created), link(if invite), id(if created), error

	userEmail := req.Email
	if userEmail == "" {
		userEmail = req.Username
	}

	if req.Password == "" && userEmail != "" {
		// Invite Flow
		token := utils.GenerateRandomToken(32)

		existingUser, err := s.UserRepo.GetUserByEmail(userEmail)

		if err == nil {
			// Found
			if existingUser.Password != nil {
				return "", "", 0, fmt.Errorf("User already registered")
			}
			// Update
			s.UserRepo.UpdateInvitedUser(existingUser.ID, req.RoleID, token)
		} else if err == sql.ErrNoRows {
			// Create
			s.UserRepo.CreateInvitedUser(userEmail, req.RoleID, token)
		} else {
			return "", "", 0, err
		}

		link := fmt.Sprintf("http://localhost:5173/set-password?token=%s", token)
		fmt.Printf("MOCK EMAIL SENT TO %s: Invite Link: %s\n", userEmail, link)
		return "invite", link, 0, nil
	}

	// Legacy Flow
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	// Need to fix CreateUser repo to accept string pass properly or struct
	// Re-mapping explicitly for repo
	hashString := string(hash)
	newUser := models.User{
		Username: req.Username,
		Password: &hashString,
		RoleID:   req.RoleID,
	}

	id, err := s.UserRepo.CreateUser(newUser)
	return "created", "", id, err
}

func (s *AdminService) UpdateUserRole(userID int, roleID int) error {
	return s.UserRepo.UpdateUserRole(userID, roleID)
}

func (s *AdminService) DeleteUser(userID string) error {
	return s.UserRepo.DeleteUser(userID)
}
