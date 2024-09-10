package auth

import (
	"context"

	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.Queries
}

func (a *AuthService) CreateUser(ctx context.Context, requestID string, dto CreateUserDTO) (*repository.User, error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()
	span.SetAttributes(attribute.String(middleware.REQUEST_ID, requestID))
	spanID := span.SpanContext().SpanID().String()
	util.LogInfo(span, requestID, spanID, "->Saving user to database start")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Fail to hash password")
		return nil, err
	}

	dto.Password = string(hashedPassword)
	user, err := a.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:    dto.Email,
		Password: dto.Password,
		Role:     repository.Role(dto.Role),
	})
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Fail to save user")
		return nil, err
	}
	util.LogInfo(span, requestID, spanID, "<-Saving user to database success")
	return &user, nil
}

func InitAuthService(repo *repository.Queries) *AuthService {
	return &AuthService{repo: repo}
}
