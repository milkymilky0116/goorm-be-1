package auth

import (
	"context"

	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.Queries
}

func (a *AuthService) CreateUser(requestID string, dto CreateUserDTO) (*repository.User, error) {
	log.Info().Str(middleware.REQUEST_ID, requestID).Msg("Saving user to database")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to hash password")
		return nil, err
	}

	dto.Password = string(hashedPassword)
	user, err := a.repo.CreateUser(context.Background(), repository.CreateUserParams{
		Email:    dto.Email,
		Password: dto.Password,
		Role:     repository.Role(dto.Role),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to save user")
		return nil, err
	}

	log.Info().Str(middleware.REQUEST_ID, requestID).Msg("Saving user complete")
	return &user, nil
}

func InitAuthService(repo *repository.Queries) *AuthService {
	return &AuthService{repo: repo}
}
