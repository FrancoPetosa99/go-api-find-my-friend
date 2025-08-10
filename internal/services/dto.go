package services

import (
	"go-api-find-my-friend/internal/models"
	"mime/multipart"
	"slices"
	"strings"
)

type AuthCredentials struct {
	Token        string `json:"token"`
	UserID       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserLastName string `json:"user_last_name"`
	UserEmail    string `json:"user_email"`
	UserPhone    string `json:"user_phone"`
}

type PetCreateDTO struct {
	Name             string                `form:"name" binding:"required"`
	Description      string                `form:"description"`
	Type             string                `form:"type" binding:"required"`
	Breed            string                `form:"breed" binding:"required"`
	LastSeenTime     string                `form:"last_seen_time" binding:"required"`
	LastSeenProvince string                `form:"last_seen_province" binding:"required"`
	LastSeenCity     string                `form:"last_seen_city" binding:"required"`
	Picture          *multipart.FileHeader `form:"picture"`
}

type UserCreateDTO struct {
	Name            string `json:"name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
	Phone           string `json:"phone"`
}

func (dto *UserCreateDTO) Validate(errors *map[string]string) bool {
	if dto.Name == "" {
		(*errors)["name"] = "Name is required"
	}

	if dto.LastName == "" {
		(*errors)["last_name"] = "Last name is required"
	}

	if dto.Email == "" {
		(*errors)["email"] = "Email is required"
	}

	if dto.Password == "" {
		(*errors)["password"] = "Password is required"
	}

	if dto.ConfirmPassword == "" {
		(*errors)["confirm_password"] = "Confirm password is required"
	}

	if dto.Password != dto.ConfirmPassword {
		(*errors)["confirm_password"] = "Passwords do not match"
	}

	if dto.Phone == "" {
		(*errors)["phone"] = "Phone is required"
	}

	if len(*errors) > 0 {
		return false
	}

	return true
}

func (dto *PetCreateDTO) Validate(errors *map[string]string) bool {
	if dto.Name == "" {
		(*errors)["name"] = "Name is required"
	}

	if dto.Description == "" {
		(*errors)["description"] = "Description is required"
	}

	if dto.Type == "" {
		(*errors)["type"] = "Type is required"
	}

	if !slices.Contains(models.PetTypes, dto.Type) {
		(*errors)["type"] = "Invalid type, must be one of: " + strings.Join(models.PetTypes, ", ")
	}

	if dto.Breed == "" {
		(*errors)["breed"] = "Breed is required"
	}

	_, breedExists := models.PetBreeds[dto.Type]
	if !breedExists {
		(*errors)["breed"] = "Invalid breed, must be one of: " + strings.Join(models.PetBreeds[dto.Type], ", ")
	}

	if dto.LastSeenTime == "" {
		(*errors)["last_seen_time"] = "Last seen time is required"
	}

	if dto.LastSeenProvince == "" {
		(*errors)["last_seen_province"] = "Province is required"
	}

	_, provinceExists := models.CitiesByProvince[dto.LastSeenProvince]
	if !provinceExists {
		(*errors)["last_seen_province"] = "Invalid province"
	}

	if dto.LastSeenCity == "" {
		(*errors)["last_seen_city"] = "City is required"
	}

	validCities := models.CitiesByProvince[dto.LastSeenProvince]
	if !slices.Contains(validCities, dto.LastSeenCity) {
		(*errors)["last_seen_city"] = "Invalid city"
	}

	if dto.Picture == nil {
		(*errors)["picture"] = "Picture is required"
	}

	if len(*errors) > 0 {
		return false
	}

	return true
}

type SearchPetsPaginationDTO struct {
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	Pets  *[]models.Pet `json:"pets"`
}

type PetUpdateDTO struct {
	Name             *string `json:"name,omitempty"`
	Type             *string `json:"type,omitempty"`
	Breed            *string `json:"breed,omitempty"`
	LastSeenTime     *string `json:"last_seen_time,omitempty"`
	LastSeenProvince *string `json:"last_seen_province,omitempty"`
	LastSeenCity     *string `json:"last_seen_city,omitempty"`
	PictureURL       *string `json:"picture_url,omitempty"`
	IsFound          *bool   `json:"is_found,omitempty"`
}

func (dto *PetUpdateDTO) Validate(errors *map[string]string) bool {
	if dto.Type != nil && !slices.Contains(models.PetTypes, *dto.Type) {
		(*errors)["type"] = "Invalid type, must be one of: " + strings.Join(models.PetTypes, ", ")
	}

	if dto.Type != nil && dto.Breed != nil {
		validBreeds, breedExists := models.PetBreeds[*dto.Type]
		if !breedExists {
			(*errors)["breed"] = "Invalid breed for this type"
		} else if !slices.Contains(validBreeds, *dto.Breed) {
			(*errors)["breed"] = "Invalid breed, must be one of: " + strings.Join(validBreeds, ", ")
		}
	}

	if dto.LastSeenProvince != nil {
		_, provinceExists := models.CitiesByProvince[*dto.LastSeenProvince]
		if !provinceExists {
			(*errors)["last_seen_province"] = "Invalid province"
		}
	}

	if dto.LastSeenProvince != nil && dto.LastSeenCity != nil {
		validCities := models.CitiesByProvince[*dto.LastSeenProvince]
		if !slices.Contains(validCities, *dto.LastSeenCity) {
			(*errors)["last_seen_city"] = "Invalid city for this province"
		}
	}

	return len(*errors) == 0
}
