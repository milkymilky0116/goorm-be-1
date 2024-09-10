package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type IdType string

const REQUEST_ID = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, IdType(REQUEST_ID), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
