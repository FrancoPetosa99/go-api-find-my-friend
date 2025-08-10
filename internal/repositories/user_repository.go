package repositories

import (
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/database"
	"go-api-find-my-friend/pkg/errors"
	"sync"

	"gorm.io/gorm"
)

type UserRepositorySQLServer struct {
	db *gorm.DB
}

var (
	userRepositoryInstance *UserRepositorySQLServer
	userRepositoryOnce     sync.Once
)

func NewUserRepositorySQLServer() *UserRepositorySQLServer {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = &UserRepositorySQLServer{
			db: database.DB,
		}
	})
	return userRepositoryInstance
}

func (r *UserRepositorySQLServer) Create(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return errors.NewInternalServerError("Failed to create user")
	}
	return nil
}

func (r *UserRepositorySQLServer) ExistsByEmail(email string) (bool, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, errors.NewInternalServerError("Failed to check user existence")
	}
	return true, nil
}

func (r *UserRepositorySQLServer) ExistsByID(id int) (bool, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, errors.NewInternalServerError("Failed to check user existence")
	}
	return true, nil
}

func (r *UserRepositorySQLServer) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get user")
	}
	return &user, nil
}

func (r *UserRepositorySQLServer) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get user by email")
	}
	return &user, nil
}
