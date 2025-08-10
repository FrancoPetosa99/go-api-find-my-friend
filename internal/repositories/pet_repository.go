package repositories

import (
	"fmt"
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/database"
	"go-api-find-my-friend/pkg/errors"
	"go-api-find-my-friend/pkg/pagination"
	"go-api-find-my-friend/pkg/storage_provider"
	"mime/multipart"
	"sync"

	"gorm.io/gorm"
)

type PetRepositorySQLServer struct {
	db              *gorm.DB
	storageProvider storage_provider.StorageProvider
}

var (
	petRepositoryInstance *PetRepositorySQLServer
	petRepositoryOnce     sync.Once
)

func NewPetRepositorySQLServer() *PetRepositorySQLServer {
	petRepositoryOnce.Do(func() {
		petRepositoryInstance = &PetRepositorySQLServer{
			db:              database.DB,
			storageProvider: storage_provider.NewCloudinary(),
		}
	})
	return petRepositoryInstance
}

func (r *PetRepositorySQLServer) Create(pet *models.Pet, picture *multipart.FileHeader) error {
	orchestrator := NewSagaOrchestrator()

	uploadPictureStep := NewUploadPictureStep(pet, picture, r.storageProvider)
	createStep := NewCreatePetStep(pet, r.db)

	orchestrator.AddSteps(uploadPictureStep, createStep)

	if err := orchestrator.Run(); err != nil {
		return err
	}

	return nil
}

func (r *PetRepositorySQLServer) GetByID(id int) (*models.Pet, error) {
	var pet models.Pet

	err := r.db.Preload("User").Where("id = ?", id).First(&pet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("pet with id %d not found", id))
		}
		return nil, errors.NewInternalServerError("An error occurred while getting pet from database")
	}

	return &pet, nil
}

func (r *PetRepositorySQLServer) Search(filter *pagination.FilterPet, search *pagination.PaginationParams) (*pagination.PaginationResult, error) {
	query := r.db.Model(&models.PetSearchResult{})

	if filter != nil {
		if filter.UserID != nil {
			query = query.Where("user_id = ?", *filter.UserID)
		}
		if filter.Type != nil {
			query = query.Where("type = ?", *filter.Type)
		}
		if filter.Breed != nil {
			query = query.Where("breed COLLATE SQL_Latin1_General_CP1_CI_AS LIKE ?", "%"+*filter.Breed+"%")
		}
		if filter.LastSeenPlace != nil {
			query = query.Where("last_seen_place COLLATE SQL_Latin1_General_CP1_CI_AS LIKE ?", "%"+*filter.LastSeenPlace+"%")
		}
	}

	var total int64
	query.Count(&total)
	if total == 0 {
		return &pagination.PaginationResult{
			Total: 0,
			Page:  search.Page,
			Size:  search.Size,
			Data:  nil,
		}, nil
	}

	sortBy := "created_at"
	if search.SortDir != "ASC" && search.SortDir != "DESC" {
		search.SortDir = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, search.SortDir))

	if search.Page > 0 && search.Size > 0 {
		offset := (search.Page - 1) * search.Size
		query = query.Offset(offset).Limit(search.Size)
	}

	pets := make([]models.PetSearchResult, 0, search.Size)

	err := query.Find(&pets).Error
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to search pets")
	}

	return &pagination.PaginationResult{
		Total: total,
		Page:  search.Page,
		Size:  search.Size,
		Data:  &pets,
	}, nil
}

func (r *PetRepositorySQLServer) Update(id int, updates map[string]interface{}) error {
	err := r.db.Model(&models.Pet{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return errors.NewInternalServerError("Failed to update pet")
	}
	return nil
}

func (r *PetRepositorySQLServer) Delete(pet *models.Pet) error {
	err := r.storageProvider.Delete(pet.PictureURL)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete pet picture")
	}

	err = r.db.Delete(&models.Pet{}, pet.ID).Error
	if err != nil {
		return errors.NewInternalServerError("Failed to delete pet")
	}

	return nil
}
