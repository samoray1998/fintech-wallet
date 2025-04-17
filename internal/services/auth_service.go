package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/samoray1998/fintech-wallet/internal/models"
	"github.com/samoray1998/fintech-wallet/internal/repositories"
)

type AuthService struct {
	UserRepo     repositories.UserRepository
	JWT_SECRET   string
	AccessExpiry time.Duration
}

func NewAuthService(repo repositories.UserRepository, jwtSecret string, accessExpiry time.Duration) *AuthService {
	return &AuthService{
		UserRepo:     repo,
		JWT_SECRET:   jwtSecret,
		AccessExpiry: accessExpiry,
	}
}

func (s *AuthService) GenerateTokens(user *models.User) (string, error) {

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.AccessExpiry).Unix(),
		"kyc":     user.KYCStatus,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWT_SECRET))
}

func (s *AuthService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
