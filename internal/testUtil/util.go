package testutil

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/milkymilky0116/goorm-be-1/internal/api"
)

func StartTestServer(t *testing.T, ctx context.Context, wg *sync.WaitGroup) net.Listener {
	listener, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer wg.Done()
		_, err := api.Run(ctx, listener)
		if err != nil {
			t.Errorf("Fail to run server: %+v", err)
			return
		}
	}()
	return listener
}
