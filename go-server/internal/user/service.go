package user

import (
	"errors"
	"fmt"
	"server/internal/config"
	"server/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *Repository
}

func (s *Service) GetAll() ([]map[string]interface{}, error) {
	return s.Repo.GetAll()
}

func (s *Service) CreateOrInvite(req CreateUserRequest) (string, string, int, error) {
	// Check if user already exists
	exists, _ := s.Repo.UsernameExists(req.Email)
	if exists {
		return "", "", 0, errors.New("User already registered")
	}

	if req.Password != "" {
		// Direct user creation with password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return "", "", 0, err
		}
		hashedPassword := string(hash)
		id, err := s.Repo.Create(req.Email, hashedPassword, req.RoleID, &req.Email, nil)
		if err != nil {
			return "", "", 0, err
		}
		// Update to Active status
		config.DB.Exec("UPDATE users SET status = 'Active' WHERE id = $1", id)
		return "created", "", id, nil
	} else {
		// Invitation flow
		token := utils.GenerateRandomToken(16)
		id, err := s.Repo.Create(req.Email, "", req.RoleID, &req.Email, &token)
		if err != nil {
			return "", "", 0, err
		}

		// Send actual email
		link := fmt.Sprintf("http://localhost:5173/set-password?token=%s", token)

		err = utils.SendInvitationEmail(req.Email, link)
		if err != nil {
			// Log error but don't fail - return link for fallback
			fmt.Printf("❌ Email send failed: %v\n", err)
			fmt.Printf("📧 Fallback - Invitation Link:\nTo: %s\nLink: %s\n", req.Email, link)
			return "invite", link, id, fmt.Errorf("invitation created but email failed: %v", err)
		}

		fmt.Printf("✅ Invitation email sent to: %s\n", req.Email)
		return "invite", link, id, nil
	}
}

func (s *Service) UpdateRole(userID, roleID int) error {
	return s.Repo.UpdateRole(userID, roleID)
}

func (s *Service) Delete(userID string) error {
	return s.Repo.Delete(userID)
}
