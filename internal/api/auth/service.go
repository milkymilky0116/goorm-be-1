package auth

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milkymilky0116/goorm-be-1/internal/api/jwt"
	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.Queries
	jwtService *jwt.JWTService
	db         *pgxpool.Pool
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

func (a *AuthService) Signin(ctx context.Context, requestID string, dto SigninDTO) (*SigninResultDTO, error) {
	ctx, span := tracer.Start(ctx, "SigninUser")
	defer span.End()
	span.SetAttributes(attribute.String(middleware.REQUEST_ID, requestID))
	spanID := span.SpanContext().SpanID().String()
	util.LogInfo(span, requestID, spanID, "->Signin user transaction start")
	tx, err := a.db.Begin(ctx)
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Fail to start transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := a.repo.WithTx(tx)
	user, err := qtx.GetUser(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.LogError(span, err, requestID, spanID, "Fail to find user by email")
			return nil, ErrUserNotFound
		}
		util.LogError(span, err, requestID, spanID, "Fail to find user by email")
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Password is incorrect")
		return nil, ErrInvalidPassword
	}
	accessToken, err := a.jwtService.GetToken("id", fmt.Sprintf("%d", user.ID), jwt.ACCESSTOKEN_EXPIRATION_TIME)
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Fail to generate AccessToken")
		return nil, ErrGenerateTokenFail
	}

	refreshToken, err := a.jwtService.GetToken("id", fmt.Sprintf("%d", user.ID), jwt.REFRESHTOKEN_EXPIRATION_TIME)
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Fail to generate RefreshToken")
		return nil, ErrGenerateTokenFail
	}
	_, err = a.repo.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:       &user.ID,
		RefreshToken: *refreshToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(jwt.REFRESHTOKEN_EXPIRATION_TIME),
			Valid: true,
		},
	})
	if err != nil {
		util.LogError(span, err, requestID, spanID, "Fail to save refresh token on db")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		util.LogError(span, err, requestID, spanID, "DB Transaction Fail")
		return nil, err
	}
	util.LogInfo(span, requestID, spanID, "<-Signin user transaction success")
	return &SigninResultDTO{AccessToken: *accessToken, RefreshToken: *refreshToken}, nil
}

func InitAuthService(repo *repository.Queries, db *pgxpool.Pool, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) *AuthService {
	jwtService := jwt.InitJWTService(publicKey, privateKey)
	return &AuthService{repo: repo, jwtService: jwtService, db: db}
}
