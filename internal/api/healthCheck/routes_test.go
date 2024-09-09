package healthcheck_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	testutil "github.com/milkymilky0116/goorm-be-1/internal/testUtil"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := testutil.StartTestServer(t, ctx, &wg)
	wg.Wait()

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/health_check", app.Listener.Addr()), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
}
