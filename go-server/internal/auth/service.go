package auth

import (
	"errors"
	"server/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *Repository
}

func (s *Service) Login(username, password string) (*LoginResponse, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Password == nil {
		return nil, errors.New("password not set. Please use invitation link")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  map[string]interface{}{"id": user.ID, "username": user.Username, "role_id": user.RoleID},
	}, nil
}

func (s *Service) SetPassword(token, password string) error {
	user, err := s.Repo.GetUserByToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	if user.Status == "Declined" {
		return errors.New("invitation was declined")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	return s.Repo.SetPassword(user.ID, string(hashedPassword))
}

func (s *Service) DeclineInvitation(token string) error {
	user, err := s.Repo.GetUserByToken(token)
	if err != nil {
		return errors.New("invalid token")
	}

	return s.Repo.DeclineInvitation(user.ID)
}

func (s *Service) GetPermissions(roleID int) ([]Permission, error) {
	return s.Repo.GetPermissionsByRoleID(roleID)
}

func (s *Service) GetFieldPermissions(roleID int) ([]FieldPermission, error) {
	return s.Repo.GetFieldPermissionsByRoleID(roleID)
}

func (s *Service) GetAllUsers() ([]map[string]interface{}, error) {
	return s.Repo.GetAllUsersSimple()
}

func (s *Service) GetEmployees() ([]map[string]interface{}, error) {
	return s.Repo.GetEmployeeNames()
}
