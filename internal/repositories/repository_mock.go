package repositories

import (
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/pagination"
	"mime/multipart"
)

type PetRepositoryMock struct {
	CreateFunc  func(pet *models.Pet, picture *multipart.FileHeader) error
	GetByIDFunc func(id int) (*models.Pet, error)
	SearchFunc  func(filter *pagination.FilterPet, search *pagination.PaginationParams) (*pagination.PaginationResult, error)
	UpdateFunc  func(id int, updates map[string]interface{}) error
	DeleteFunc  func(id int) error
}

func (m *PetRepositoryMock) Create(pet *models.Pet, picture *multipart.FileHeader) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(pet, picture)
	}
	return nil
}

func (m *PetRepositoryMock) GetByID(id int) (*models.Pet, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *PetRepositoryMock) Search(filter *pagination.FilterPet, search *pagination.PaginationParams) (*pagination.PaginationResult, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(filter, search)
	}
	return nil, nil
}

func (m *PetRepositoryMock) Update(id int, updates map[string]interface{}) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, updates)
	}
	return nil
}

func (m *PetRepositoryMock) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

type UserRepositoryMock struct {
	CreateFunc        func(user *models.User) error
	GetByIDFunc       func(id int) (*models.User, error)
	GetByEmailFunc    func(email string) (*models.User, error)
	ExistsByEmailFunc func(email string) (bool, error)
	ExistsByIDFunc    func(id int) (bool, error)
}

func (m *UserRepositoryMock) Create(user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *UserRepositoryMock) GetByEmail(email string) (*models.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, nil
}

func (m *UserRepositoryMock) ExistsByEmail(email string) (bool, error) {
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(email)
	}
	return false, nil
}

func (m *UserRepositoryMock) ExistsByID(id int) (bool, error) {
	if m.ExistsByIDFunc != nil {
		return m.ExistsByIDFunc(id)
	}
	return false, nil
}

func (m *UserRepositoryMock) GetByID(id int) (*models.User, error) {
	return nil, nil
}
