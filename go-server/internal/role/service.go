package role

type Service struct {
	Repo *Repository
}

func (s *Service) GetAll() ([]Role, error) {
	return s.Repo.GetAll()
}

func (s *Service) Create(name string) (int, error) {
	return s.Repo.Create(name)
}

func (s *Service) Delete(id int) error {
	return s.Repo.Delete(id)
}

func (s *Service) GetPermissions(roleID string) ([]Permission, error) {
	return s.Repo.GetPermissions(roleID)
}

func (s *Service) AddOrUpdatePermission(roleID int, resource, action string) (string, int, error) {
	return s.Repo.AddOrUpdatePermission(roleID, resource, action)
}

func (s *Service) DeletePermission(roleID, resource, action string) error {
	return s.Repo.DeletePermission(roleID, resource, action)
}

func (s *Service) GetFieldPermissions(roleID int) ([]FieldPermission, error) {
	return s.Repo.GetFieldPermissions(roleID)
}

func (s *Service) UpdateFieldPermission(roleID int, resource, field string, canView, canEdit bool) error {
	return s.Repo.UpdateFieldPermission(roleID, resource, field, canView, canEdit)
}
