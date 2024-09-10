package auth

import (
	"context"
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	repo      *repository.Queries
	validator *validator.Validate
}

func (app *AuthController) SignupController(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserDTO
	err := util.ReadBody(r, &dto)
	if err != nil {
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	err = app.validator.Struct(dto)
	if err != nil {
		util.HandleValidatorError(w, err, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	dto.Password = string(hashedPassword)
	user, err := app.repo.CreateUser(context.Background(), repository.CreateUserParams{
		Email:    dto.Email,
		Password: dto.Password,
		Role:     repository.Role(dto.Role),
	})
	if err != nil {
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	err = util.WriteJson(w, user, http.StatusOK)
	if err != nil {
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func InitAuthController(repo *repository.Queries, validate *validator.Validate) *AuthController {
	return &AuthController{
		repo:      repo,
		validator: validate,
	}
}
