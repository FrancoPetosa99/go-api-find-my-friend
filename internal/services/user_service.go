package services

import (
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/internal/repositories"
	"go-api-find-my-friend/pkg/errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository repositories.UserRepository
}

var (
	userServiceInstance *UserService
	userServiceOnce     sync.Once
)

func NewUserService() *UserService {
	userServiceOnce.Do(func() {
		userServiceInstance = &UserService{
			userRepository: repositories.NewUserRepository(),
		}
	})
	return userServiceInstance
}

func (s *UserService) CreateUser(dto *UserCreateDTO) (*models.User, error) {
	err := s.CheckForDuplicates(dto.Email)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     dto.Name,
		LastName: dto.LastName,
		Email:    dto.Email,
		Password: hashedPassword,
		Phone:    dto.Phone,
	}

	err = s.userRepository.Create(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalServerError("An error occurred while hashing password")
	}
	return string(hashedPassword), nil
}

func (s *UserService) CheckUserExists(ID int) (bool, error) {
	exists, err := s.userRepository.ExistsByID(ID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *UserService) CheckForDuplicates(email string) error {
	exists, err := s.userRepository.ExistsByEmail(email)
	if err != nil {
		return err
	}

	if exists {
		return errors.NewConflictError("Already exists user with email " + email)
	}

	return nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
