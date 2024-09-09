package auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (a *AuthController) SignupController(w http.ResponseWriter, r *http.Request) {
	var result map[string]string
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(result["password"]), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result["password"] = string(hashedPassword)
	jsonBody, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func InitAuthController() *AuthController {
	return &AuthController{}
}
