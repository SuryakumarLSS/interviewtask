package resource

import "errors"

type Service struct {
	Repo *Repository
}

func (s *Service) GetAll(resource string, roleID int) ([]map[string]interface{}, error) {
	// Check permission
	allowed, _ := s.Repo.HasPermission(roleID, resource, "read")
	if !allowed {
		return nil, errors.New("permission denied")
	}

	return s.Repo.GetAll(resource, roleID)
}

func (s *Service) Create(resource string, data map[string]interface{}, roleID int) (int, error) {
	// Check permission
	allowed, _ := s.Repo.HasPermission(roleID, resource, "create")
	if !allowed {
		return 0, errors.New("permission denied")
	}

	return s.Repo.Create(resource, data, roleID)
}

func (s *Service) Update(resource, id string, data map[string]interface{}, roleID int) error {
	// Check permission
	allowed, _ := s.Repo.HasPermission(roleID, resource, "update")
	if !allowed {
		return errors.New("permission denied")
	}

	return s.Repo.Update(resource, id, data, roleID)
}

func (s *Service) Delete(resource, id string, roleID int) error {
	// Check permission
	allowed, _ := s.Repo.HasPermission(roleID, resource, "delete")
	if !allowed {
		return errors.New("permission denied")
	}

	return s.Repo.Delete(resource, id)
}
