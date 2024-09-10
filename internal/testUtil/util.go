package testutil

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/milkymilky0116/goorm-be-1/internal/api"
	"github.com/milkymilky0116/goorm-be-1/internal/configuration"
	"github.com/milkymilky0116/goorm-be-1/internal/db/repository"
	"github.com/pressly/goose/v3"
)

type TestApp struct {
	Repo     *repository.Queries
	Listener net.Listener
}

const MIGRATION_DIR = "migration"

func MigrateDB(connectionString string) error {
	db, err := sql.Open("postgres", connectionString+"sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.Up(db, "migration"); err != nil {
		return err
	}
	return nil
}

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
		t.Errorf("Fail to set root directory: %v", err)
	}

	config, err := configuration.InitConfiguration()
	if err != nil {
		t.Errorf("Fail to load configuration: %v", err)
	}
	connectionString, err := LaunchPostgresContainer(&config.Database)
	if err != nil {
		t.Errorf("Fail to launch database container: %v", err)
	}
	conn, err := pgxpool.New(context.Background(), *connectionString)
	if err != nil {
		t.Errorf("Fail to connect to db: %v", err)
	}
	err = MigrateDB(*connectionString)
	if err != nil {
		t.Errorf("Fail to migrate db: %v", err)
	}
	repo := repository.New(conn)
	listener, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		_, err := api.Run(ctx, listener, conn)
		if err != nil {
			t.Errorf("Fail to run server: %+v", err)
			return
		}
	}()
	return TestApp{
		Repo:     repo,
		Listener: listener,
	}
}
