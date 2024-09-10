package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrBody struct {
	Msg string `json:"msg"`
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
