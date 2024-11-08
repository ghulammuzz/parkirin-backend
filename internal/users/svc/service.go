package svc

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	userEntity "github.com/ghulammuzz/backend-parkerin/internal/users/entity"
	userRepo "github.com/ghulammuzz/backend-parkerin/internal/users/repo"
)

type UserService interface {
	ListUser(page, limit int) (*userEntity.UserListResponse, error)
	RegisterUser(user *userEntity.UserRegisterRequest) error
	LoginUser(user *userEntity.UserLoginRequest) (string, error)
	LoginStore(user *userEntity.UserLoginRequest) (string, error)
	GetUserDetails(userID int) (*userEntity.UserDetailResponse, error)
	IsPhoneNumberExists(phone string) (bool, error)
}

type userService struct {
	userRepo userRepo.UserRepository
}

func (s *userService) ListUser(page int, limit int) (*userEntity.UserListResponse, error) {
	users, err := s.userRepo.List(page, limit)
	if err != nil {
		return &userEntity.UserListResponse{}, err
	}
	return users, nil
}

func (s *userService) LoginStore(user *userEntity.UserLoginRequest) (string, error) {
	dbUser, err := s.userRepo.LoginUser(user)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"apps":    "parkirin-backend",
		"user_id": dbUser.ID,
		"role":    dbUser.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set in environment variables")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) IsPhoneNumberExists(phone string) (bool, error) {
	return s.userRepo.IsPhoneNumberExists(phone)
}

func (s *userService) GetUserDetails(userID int) (*userEntity.UserDetailResponse, error) {
	return s.userRepo.Detail(userID)
}

func (s *userService) RegisterUser(user *userEntity.UserRegisterRequest) error {
	if user.Role != "tukang" && user.Role != "store" {
		return errors.New("invalid role")
	}
	if user.Role == "tukang" {
		user.IsVerified = true
	}

	return s.userRepo.Create(user)
}

func (s *userService) LoginUser(user *userEntity.UserLoginRequest) (string, error) {

	dbUser, err := s.userRepo.LoginUser(user)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"apps":    "parkirin-backend",
		"user_id": dbUser.ID,
		"role":    dbUser.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set in environment variables")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewUserService(userRepo userRepo.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
