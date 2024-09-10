package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/milkymilky0116/goorm-be-1/internal/api/middleware"
	"github.com/rs/zerolog/log"
)

type ErrBody struct {
	Msg string `json:"msg"`
}

func GetRequestID(w http.ResponseWriter, r *http.Request) *string {
	requestID, ok := r.Context().Value(middleware.IdType(middleware.REQUEST_ID)).(string)
	if !ok {
		log.Error().Msg("requestID was not found in the context")
		HandleError(w, errors.New("fail to find request id"), http.StatusInternalServerError)
		return nil
	}
	return &requestID
}

func ReadBody(r *http.Request, data any) error {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, data any, statusCode int) error {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}
	w.Write(jsonBody)
	return nil
}

func HandleError(w http.ResponseWriter, err error, statusCode int) {
	WriteJson(w, &ErrBody{Msg: err.Error()}, statusCode)
}

func HandleValidatorError(w http.ResponseWriter, err error, statusCode int) {
	var msg string
	for _, err := range err.(validator.ValidationErrors) {
		errMsg := fmt.Sprintf("Field validation failed on %s\n", err.Field())
		msg += errMsg
	}
	WriteJson(w, &ErrBody{Msg: msg}, statusCode)
}
