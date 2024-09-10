package auth

import (
	"net/http"

	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/go-playground/validator/v10"
)

var tracer = otel.Tracer("SignUp")

type AuthController struct {
	authService *AuthService
	validate    *validator.Validate
}

func (app *AuthController) SignupController(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "SignUp")
	defer span.End()
	requestID := util.GetRequestID(w, r)
	spanID := span.SpanContext().SpanID().String()

	span.SetAttributes(attribute.String(middleware.REQUEST_ID, *requestID))
	util.LogInfo(span, *requestID, spanID, "Starting register user")
	var dto CreateUserDTO
	err := util.ReadBody(r, &dto)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusInternalServerError, "Fail to read body")
		return
	}
	err = app.validate.Struct(dto)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusInternalServerError, "Fail to validate body")
		return
	}
	user, err := app.authService.CreateUser(ctx, *requestID, dto)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusInternalServerError, "Fail to create user")
		return
	}

	err = util.WriteJson(w, user, http.StatusOK)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusInternalServerError, "Fail to write json")
		return
	}
	util.LogInfo(span, *requestID, spanID, "Finish registering user")
}

func InitAuthController(repo *repository.Queries, validate *validator.Validate) *AuthController {
	authService := InitAuthService(repo)
	return &AuthController{
		authService: authService,
		validate:    validate,
	}
}
