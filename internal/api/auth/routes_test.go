package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	testutil "github.com/milkymilky0116/goorm-be-1/internal/testUtil"
	"github.com/stretchr/testify/assert"
)

func TestAuthController(t *testing.T) {
	t.Run("signup should return 200 when valid input was given", func(t *testing.T) {
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		listener := testutil.StartTestServer(t, ctx, &wg)
		wg.Wait()

		client := &http.Client{}
		body := map[string]string{
			"email":    "test@naver.com",
			"password": "123",
			"role":     "STUDENT",
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Fail to marshal json body : %+v", err)
		}
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/auth/signup", listener.Addr()), bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Errorf("Fail to request url: %+v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Fail to request url: %+v", err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var result map[string]string
		if err != nil {
			t.Errorf("Fail to read body: %+v", err)
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Errorf("Fail to decode json: %+v", err)
		}

		assert.Equal(t, "test@naver.com", result["email"])
		assert.NotEqual(t, "123", result["password"])
		assert.Equal(t, "STUDENT", result["role"])

	})
}
