package services

import (
	"fmt"
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/internal/repositories"
	"go-api-find-my-friend/pkg/errors"
	"go-api-find-my-friend/pkg/pagination"
	"sync"
	"time"
)

type PetService struct {
	petRepository repositories.PetRepository
	fileService   *FileService
}

var (
	petServiceInstance *PetService
	petServiceOnce     sync.Once
)

func NewPetService() *PetService {
	petServiceOnce.Do(func() {
		petServiceInstance = &PetService{
			petRepository: repositories.NewPetRepository(),
			fileService:   NewFileService(),
		}
	})
	return petServiceInstance
}

func (s *PetService) CreatePet(dto *PetCreateDTO, userID int) (*models.Pet, error) {
	lastSeenTime, err := time.Parse("2006-02-01", dto.LastSeenTime)
	if err != nil {
		fmt.Println("Error parsing last seen time:", err)
		fmt.Println("Error parsing last seen time:", dto.LastSeenTime)
		return nil, errors.NewBadRequestError("Invalid date format. Expected format: dd-mm-yyyy")
	}

	pet := models.Pet{
		Name:          dto.Name,
		Description:   dto.Description,
		Type:          dto.Type,
		Breed:         dto.Breed,
		UserID:        userID,
		LastSeenTime:  lastSeenTime,
		LastSeenPlace: dto.LastSeenProvince + ", " + dto.LastSeenCity,
		IsFound:       false,
	}

	err = s.petRepository.Create(&pet, dto.Picture)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (s *PetService) SearchPets(filters *pagination.FilterPet, paginationParams *pagination.PaginationParams) (*pagination.PaginationResult, error) {
	customConfig := pagination.PaginationConfig{
		DefaultPage:    1,
		DefaultSize:    10,
		MaxSize:        50,
		DefaultSortBy:  "created_at",
		DefaultSortDir: "DESC",
	}
	pagination.NormalizeParams(paginationParams, customConfig)

	result, err := s.petRepository.Search(filters, paginationParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PetService) GetPetByID(petID int) (*models.Pet, error) {
	pet, err := s.petRepository.GetByID(petID)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *PetService) UpdatePet(userID int, petID int, dto *PetUpdateDTO) error {
	pet, err := s.GetPetByID(petID)
	if err != nil {
		return err
	}

	if pet.UserID != userID {
		return errors.NewForbiddenError("You can only update your own pets")
	}

	updates := make(map[string]interface{})

	if dto.Name != nil {
		updates["name"] = *dto.Name
	}
	if dto.Type != nil {
		updates["type"] = *dto.Type
	}
	if dto.Breed != nil {
		updates["breed"] = *dto.Breed
	}
	if dto.LastSeenTime != nil {
		lastSeenTime, err := time.Parse("02-01-2006", *dto.LastSeenTime)
		if err != nil {
			return errors.NewBadRequestError("Invalid date format. Expected format: dd-mm-yyyy")
		}
		updates["last_seen_time"] = lastSeenTime
	}
	if dto.LastSeenProvince != nil && dto.LastSeenCity != nil {
		updates["last_seen_place"] = *dto.LastSeenProvince + ", " + *dto.LastSeenCity
	}
	if dto.PictureURL != nil {
		updates["picture_url"] = *dto.PictureURL
	}
	if dto.IsFound != nil {
		updates["is_found"] = *dto.IsFound
	}

	err = s.petRepository.Update(petID, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *PetService) UpdatePetAsFound(userID int, petID int) error {
	pet, err := s.GetPetByID(petID)
	if err != nil {
		return err
	}

	if pet.UserID != userID {
		return errors.NewForbiddenError("You can only update your own pets")
	}

	if pet.IsFound {
		return errors.NewBadRequestError("Pet already marked as found")
	}

	updates := make(map[string]interface{})
	updates["is_found"] = true

	err = s.petRepository.Update(petID, updates)
	if err != nil {
		return err
	}

	return nil
}

func (s *PetService) DeletePet(userID int, petID int) error {
	pet, err := s.GetPetByID(petID)
	if err != nil {
		return err
	}

	if pet.UserID != userID {
		return errors.NewForbiddenError("You can only delete your own pets")
	}

	err = s.petRepository.Delete(pet)
	if err != nil {
		return err
	}

	return nil
}
