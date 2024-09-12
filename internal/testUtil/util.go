package testutil

import (
	"context"
	"crypto/ed25519"
	"errors"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/milkymilky0116/goorm-be-1/internal/api"
	"github.com/milkymilky0116/goorm-be-1/internal/api/jwt"
	"github.com/milkymilky0116/goorm-be-1/internal/configuration"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/milkymilky0116/goorm-be-1/internal/tracing"
	"github.com/milkymilky0116/goorm-be-1/internal/util"
)

type TestApp struct {
	Repo       *repository.Queries
	Listener   net.Listener
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

const MIGRATION_DIR = "migration"

func SetRootDirectory() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			break
		}

		wd = filepath.Dir(wd)
		if wd == "/" {
			return errors.New("cannot find root directory")
		}
	}

	os.Chdir(wd)
	return nil
}

func StartTestServer(t *testing.T, ctx context.Context, wg *sync.WaitGroup) TestApp {
	err := SetRootDirectory()
	if err != nil {
		t.Fatalf("Fail to set root directory: %v", err)
	}

	config, err := configuration.InitConfiguration()
	if err != nil {
		t.Fatalf("Fail to load configuration: %v", err)
	}

	tracingProvider, err := tracing.InitTracing("goorm-class-test", config.Jaeger.Host, config.Jaeger.Port)
	defer func() {
		if err := tracingProvider.Shutdown(context.Background()); err != nil {
			t.Fatalf("Fail to shutting down tracing provider")
		}
	}()
	tracingProvider.Tracer("goorm-class-test")

	if err != nil {
		t.Fatalf("Fail to set tracing provider")
	}

	connectionString, err := LaunchPostgresContainer(&config.Database)
	if err != nil {
		t.Fatalf("Fail to launch database container: %v", err)
	}
	conn, err := pgxpool.New(context.Background(), *connectionString)
	if err != nil {
		t.Fatalf("Fail to connect to db: %v", err)
	}
	err = util.MigrateDB(*connectionString)
	if err != nil {
		t.Fatalf("Fail to migrate db: %v", err)
	}
	repo := repository.New(conn)
	listener, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Fatal(err)
	}
	publicKey, privateKey, err := jwt.InitKey()
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	go func() {
		var err error
		_, err = api.Run(ctx, listener, conn, publicKey, privateKey)
		if err != nil {
			t.Errorf("Fail to run server: %+v", err)
			return
		}
	}()
	return TestApp{
		Repo:       repo,
		Listener:   listener,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}
