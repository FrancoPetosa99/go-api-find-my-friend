package controllers

import (
	"time"
)

type PetCreateResponse struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Breed         string    `json:"breed"`
	LastSeenTime  time.Time `json:"last_seen_time"`
	LastSeenPlace string    `json:"last_seen_place"`
}

type PetDetailDTO struct {
	PetID         int       `json:"pet_id"`
	OwnerID       int       `json:"owner_id"`
	OwnerName     string    `json:"owner_name"`
	OwnerLastName string    `json:"owner_last_name"`
	OwnerEmail    string    `json:"owner_email"`
	OwnerPhone    string    `json:"owner_phone"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	Breed         string    `json:"breed"`
	LastSeenTime  time.Time `json:"last_seen_time"`
	LastSeenPlace string    `json:"last_seen_place"`
	PictureURL    string    `json:"picture_url"`
	IsFound       bool      `json:"is_found"`
	CanEdit       bool      `json:"can_edit"`
	CanDelete     bool      `json:"can_delete"`
}

type SearchPetsPaginationDTO struct {
	Page          int    `json:"page" form:"page"`
	Size          int    `json:"size" form:"size"`
	SortDir       string `json:"sort_dir" form:"sort_dir"`
	Type          string `json:"type" form:"type"`
	Breed         string `json:"breed" form:"breed"`
	LastSeenPlace string `json:"last_seen_place" form:"last_seen_place"`
}

type UserCreateResponse struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
}

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
