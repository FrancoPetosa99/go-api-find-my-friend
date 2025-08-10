package models

import (
	"time"
)

type Pet struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null;default:''"`
	Type          string    `json:"type" gorm:"not null"`
	Breed         string    `json:"breed"`
	UserID        int       `json:"user_id" gorm:"type:int;not null;constraint:OnDelete:CASCADE"`
	User          User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	LastSeenTime  time.Time `json:"last_seen_time" gorm:"not null"`
	LastSeenPlace string    `json:"last_seen_place" gorm:"not null"`
	IsFound       bool      `json:"is_found" gorm:"default:false"`
	PictureURL    string    `json:"picture_url"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PetSearchResult struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null;default:''"`
	Type          string    `json:"type" gorm:"not null"`
	Breed         string    `json:"breed"`
	UserID        int       `json:"user_id" gorm:"type:int;not null"`
	LastSeenTime  time.Time `json:"last_seen_time" gorm:"not null"`
	LastSeenPlace string    `json:"last_seen_place" gorm:"not null"`
	IsFound       bool      `json:"is_found" gorm:"default:false"`
	PictureURL    string    `json:"picture_url"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PetSearchResult) TableName() string {
	return "pets"
}
