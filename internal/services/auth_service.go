package services

import (
	"go-api-find-my-friend/pkg/config"
	"go-api-find-my-friend/pkg/errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.NewUnauthorizedError("invalid credentials")
)

var (
	authServiceInstance *AuthService
	authServiceOnce     sync.Once
)

type AuthService struct {
	secretKey   string
	userService *UserService
}

func NewAuthService() *AuthService {
	authServiceOnce.Do(func() {
		authServiceInstance = &AuthService{
			secretKey:   config.ConfigInstance.JWT.Secret,
			userService: NewUserService(),
		}
	})
	return authServiceInstance
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(userID int, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) AuthenticateUser(email string, password string) (string, error) {
	user, err := s.userService.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := checkPassword(password, user.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func checkPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
