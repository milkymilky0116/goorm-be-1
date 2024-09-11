package auth

import (
	"crypto/ed25519"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
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
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusBadRequest, "Fail to validate body")
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

func (app *AuthController) SigninController(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "SignIn")
	defer span.End()
	requestID := util.GetRequestID(w, r)
	spanID := span.SpanContext().SpanID().String()
	span.SetAttributes(attribute.String(middleware.REQUEST_ID, *requestID))
	util.LogInfo(span, *requestID, spanID, "Starting Signin user")

	var dto SigninDTO
	err := util.ReadBody(r, &dto)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusInternalServerError, "Fail to read body")
		return
	}

	err = app.validate.Struct(dto)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusBadRequest, "Fail to validate body")
		return
	}

	tokens, err := app.authService.Signin(ctx, *requestID, dto)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusBadRequest, "Fail to Signin")
		return
	}

	err = util.WriteJson(w, tokens, http.StatusOK)
	if err != nil {
		util.HandleErrAndLog(w, span, err, *requestID, spanID, http.StatusInternalServerError, "Fail to write json")
		return
	}
	util.LogInfo(span, *requestID, spanID, "Finish Signin user")
}

func InitAuthController(repo *repository.Queries, validate *validator.Validate, db *pgxpool.Pool, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) *AuthController {
	authService := InitAuthService(repo, db, privateKey, publicKey)
	return &AuthController{
		authService: authService,
		validate:    validate,
	}
}
