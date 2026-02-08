package services

import (
	"database/sql"
	"errors"
	"server/models"
	"server/repositories"
	"server/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, errors.New("User not found")
		}
		return "", nil, err
	}

	if user.Password == nil {
		return "", nil, errors.New("Account pending activation. Please use the invitation link.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("Invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) SetPassword(token, password string) error {
	userID, err := s.UserRepo.GetUserByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Invalid or expired invitation token")
		}
		return err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return s.UserRepo.ActivateUser(userID, string(hash))
}

func (s *AuthService) DeclineInvitation(token string) error {
	userID, err := s.UserRepo.GetUserByToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Invalid or expired invitation token")
		}
		return err
	}

	return s.UserRepo.DeclineUser(userID)
}
