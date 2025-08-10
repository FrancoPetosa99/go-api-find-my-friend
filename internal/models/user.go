package models

import (
	"time"
)

type User struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"not null"`
	LastName  string     `json:"last_name" gorm:"not null"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"-" gorm:"not null"` // "-" oculta el campo en JSON
	Phone     string     `json:"phone"`
	Pets      []Pet      `json:"pets,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
