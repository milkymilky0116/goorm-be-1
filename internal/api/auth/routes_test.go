package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/milkymilky0116/goorm-be-1/internal/api/auth"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	testutil "github.com/milkymilky0116/goorm-be-1/internal/testUtil"
	"github.com/stretchr/testify/assert"
)

func TestAuthController(t *testing.T) {
	t.Run("signup should return 200 when valid input was given", func(t *testing.T) {
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		app := testutil.StartTestServer(t, ctx, &wg)

		client := &http.Client{}
		body := auth.CreateUserDTO{
			Email:    "test@naver.com",
			Password: "Abcd1234!",
			Role:     "student",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Fail to marshal json body : %+v", err)
		}
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/auth/signup", app.Listener.Addr()), bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Fail to request url: %+v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Fail to request url: %+v", err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Header.Set("Content-Type", "application/json")
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Fail to read body: %+v", err)
		}

		var result repository.User
		if err != nil {
			t.Errorf("Fail to read body: %+v", err)
		}

		err = json.Unmarshal(respBody, &result)
		if err != nil {
			t.Errorf("Fail to decode json: %+v", err)
		}

		assert.Equal(t, "test@naver.com", result.Email)
		assert.NotEqual(t, "Abcd1234!", result.Password)
		assert.Equal(t, repository.RoleStudent, result.Role)

		user, err := app.Repo.GetUser(context.Background(), result.ID)
		if err != nil {
			t.Errorf("Fail to get user : %v", err)
		}
		assert.Equal(t, "test@naver.com", user.Email)
		assert.NotEqual(t, "Abcd123!", user.Password)
		assert.Equal(t, repository.RoleStudent, user.Role)
		defer wg.Done()
	})
	t.Run("signup should return 500 when invalid input was given", func(t *testing.T) {
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		app := testutil.StartTestServer(t, ctx, &wg)

		client := &http.Client{}
		testcases := []struct {
			name string
			body map[string]string
		}{
			{
				name: "email field is not valid",
				body: map[string]string{
					"email":    "test",
					"password": "123",
					"role":     "student",
				},
			},
			{
				name: "role field is not valid",
				body: map[string]string{
					"email":    "test@naver.com",
					"password": "123",
					"role":     "test",
				},
			},
			{
				name: "email field is missing",
				body: map[string]string{
					"password": "123",
					"role":     "test",
				},
			},
			{
				name: "password field is missing",
				body: map[string]string{
					"email": "test@naver.com",
					"role":  "test",
				},
			},
			{
				name: "role field is missing",
				body: map[string]string{
					"email":    "test@naver.com",
					"password": "123",
				},
			},
			{
				name: "password field is invalid",
				body: map[string]string{
					"email":    "test@naver.com",
					"password": "123",
					"role":     "student",
				},
			},
		}
		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				jsonBody, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("Fail to parse json: %v", err)
				}
				req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/auth/signup", app.Listener.Addr()), bytes.NewBuffer(jsonBody))
				if err != nil {
					t.Fatalf("Fail to create request url: %v", err)
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatalf("Fail to request url: %v", err)
				}
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		}

		defer wg.Done()
	})
}
