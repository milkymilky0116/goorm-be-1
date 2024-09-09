package auth

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo *repository.Queries
}

func (app *AuthController) SignupController(w http.ResponseWriter, r *http.Request) {
	var dto repository.CreateUserParams
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(dto)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dto.Password = string(hashedPassword)
	user, err := app.repo.CreateUser(context.Background(), repository.CreateUserParams{
		Email:    dto.Email,
		Password: dto.Password,
		Role:     dto.Role,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonBody, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func InitAuthController(repo *repository.Queries) *AuthController {
	return &AuthController{
		repo: repo,
	}
}
