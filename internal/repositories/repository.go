package repositories

import (
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/pagination"
	"mime/multipart"
)

// PetRepository define los métodos para operaciones con mascotas
type PetRepository interface {
	Create(pet *models.Pet, picture *multipart.FileHeader) error
	GetByID(id int) (*models.Pet, error)
	Search(filter *pagination.FilterPet, search *pagination.PaginationParams) (*pagination.PaginationResult, error)
	Update(id int, updates map[string]interface{}) error
	Delete(pet *models.Pet) error
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByID(id int) (bool, error)
}

type ImageRepository interface {
	Upload(file *multipart.FileHeader) (string, error)
}

func NewPetRepository() PetRepository {
	return NewPetRepositorySQLServer()
}

func NewUserRepository() UserRepository {
	return NewUserRepositorySQLServer()
}
