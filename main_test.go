package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer wg.Done()
		_, err := run(ctx, listener)
		if err != nil {
			t.Errorf("Fail to run server: %+v", err)
			return
		}
	}()
	wg.Wait()

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/health_check", listener.Addr()), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
}
