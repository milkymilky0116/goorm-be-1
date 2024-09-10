package auth

import (
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"github.com/rs/zerolog/log"

	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService *AuthService
	validate    *validator.Validate
}

func (app *AuthController) SignupController(w http.ResponseWriter, r *http.Request) {
	requestID := util.GetRequestID(w, r)

	log.Info().Str(middleware.REQUEST_ID, *requestID).Msg("Starting registering user")
	var dto CreateUserDTO
	err := util.ReadBody(r, &dto)
	if err != nil {
		log.Err(err).Msg("Fail to read body")
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	err = app.validate.Struct(dto)
	if err != nil {
		log.Err(err).Msg("Fail to validate body")
		util.HandleValidatorError(w, err, http.StatusBadRequest)
		return
	}
	user, err := app.authService.CreateUser(*requestID, dto)
	if err != nil {
		log.Err(err).Msg("Fail to save user")
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	err = util.WriteJson(w, user, http.StatusOK)
	if err != nil {
		log.Err(err).Msg("Fail to write json")
		util.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	log.Info().Str(middleware.REQUEST_ID, *requestID).Msg("Finish registering user")
}

func InitAuthController(repo *repository.Queries, validate *validator.Validate) *AuthController {
	authService := InitAuthService(repo)
	return &AuthController{
		authService: authService,
		validate:    validate,
	}
}
