package services

import (
	"fmt"
	"server/repositories"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResourceService struct {
	ResourceRepo *repositories.ResourceRepository
	UserRepo     *repositories.UserRepository
	RoleRepo     *repositories.RoleRepository
}

var RequiredFields = map[string][]string{
	"employees": {"name", "position", "salary", "department"},
	"projects":  {"name", "assigned_to", "status", "budget"},
	"orders":    {"customer_name", "amount", "status", "order_date"},
}

func (s *ResourceService) GetAll(resource string) ([]map[string]interface{}, error) {
	return s.ResourceRepo.GetAll(resource)
}

func (s *ResourceService) Create(resource string, data map[string]interface{}) (int, error) {
	// Validation
	reqFields := RequiredFields[resource]
	var missing []string
	for _, f := range reqFields {
		if _, ok := data[f]; !ok {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		return 0, fmt.Errorf("Missing required fields: %s", strings.Join(missing, ", "))
	}

	return s.ResourceRepo.Create(resource, data)
}

func (s *ResourceService) Update(resource string, id string, data map[string]interface{}) error {
	return s.ResourceRepo.Update(resource, id, data)
}

func (s *ResourceService) Delete(resource string, id string) error {
	return s.ResourceRepo.Delete(resource, id)
}

func (s *ResourceService) GetEmployeeNames() ([]gin.H, error) {
	return s.ResourceRepo.GetEmployeeNames()
}

func (s *ResourceService) GetUserList() ([]gin.H, error) {
	return s.UserRepo.GetAllUsers()
}

func (s *ResourceService) CheckAccess(roleID int, resource string, action string) (bool, error) {
	return s.RoleRepo.CheckAccess(roleID, resource, action)
}

func (s *ResourceService) GetPermissions(roleID int) ([]gin.H, error) {
	perms, err := s.RoleRepo.GetPermissionsByIntID(roleID)
	if err != nil {
		return nil, err
	}
	// Convert to gin.H or just return models
	var result []gin.H
	for _, p := range perms {
		result = append(result, gin.H{
			"id": p.ID, "role_id": p.RoleID, "resource": p.Resource, "action": p.Action, "attributes": p.Attributes,
		})
	}
	return result, nil
}
