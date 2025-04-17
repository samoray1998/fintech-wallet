package services

import (
	"errors"
	"time"

	"github.com/samoray1998/fintech-wallet/internal/models"
	"github.com/samoray1998/fintech-wallet/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	UserRepo   repositories.UserRepository
	bcryptCost int
}

func NewUserService(repo repositories.UserRepository, bcryptCost int) *UserServices {
	return &UserServices{
		UserRepo:   repo,
		bcryptCost: bcryptCost,
	}
}

func (s *UserServices) Register(user *models.User) (*models.User, error) {

	existingUser, _ := s.UserRepo.FindByEmail(user.Email)

	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), s.bcryptCost)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	user.KYCStatus = "unverified"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.UserRepo.CreateUser(user)
}

func (s *UserServices) GetUserByID(id string) (*models.User, error) {

	user, err := s.UserRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *UserServices) VerifyCredentials(email, hashedPassword string) (*models.User, error) {

	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	/// let's compiare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil

}

func (s *UserServices) UpdateKYCStatus(userID, status string) (*models.User, error) {
	return s.UserRepo.UpdateKYCStatus(userID, status)

}


