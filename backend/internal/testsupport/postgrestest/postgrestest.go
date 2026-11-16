// Package postgrestest starts a throwaway PostgreSQL for integration tests,
// with every Liquibase migration applied by the same runner image used in
// compose and CI (database/Dockerfile). It requires Docker.
package postgrestest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	image        = "postgres:18-alpine"
	databaseName = "openboard"
	credential   = "openboard"
	networkAlias = "postgres"
)

// Database is a running throwaway PostgreSQL with the schema migrated.
type Database struct {
	URL       string
	resources []testcontainers.Container
	network   *testcontainers.DockerNetwork
}

// Start launches PostgreSQL, applies the Liquibase changelog and returns the
// connection URL. Call Stop when done.
func Start(ctx context.Context) (*Database, error) {
	db := &Database{}
	if err := db.start(ctx); err != nil {
		db.Stop()
		return nil, err
	}
	return db, nil
}

func (db *Database) start(ctx context.Context) error {
	net, err := network.New(ctx)
	if err != nil {
		return fmt.Errorf("create network: %w", err)
	}
	db.network = net

	pg, err := postgres.Run(ctx, image,
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(credential),
		postgres.WithPassword(credential),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{networkAlias}, net),
	)
	if pg != nil {
		db.resources = append(db.resources, pg)
	}
	if err != nil {
		return fmt.Errorf("start postgres: %w", err)
	}

	if err := db.migrate(ctx, net.Name); err != nil {
		return err
	}

	url, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("connection string: %w", err)
	}
	db.URL = url
	return nil
}

func (db *Database) migrate(ctx context.Context, networkName string) error {
	dir, err := databaseDir()
	if err != nil {
		return err
	}
	runner, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context: dir, Dockerfile: "Dockerfile",
				Repo: "openboard-migrations", Tag: "test", KeepImage: true,
			},
			Networks: []string{networkName},
			Cmd: []string{
				"--changelog-file=migrations/changelog/db.changelog-master.yaml",
				fmt.Sprintf("--url=jdbc:postgresql://%s:5432/%s", networkAlias, databaseName),
				"--username=" + credential,
				"--password=" + credential,
				"update",
			},
			WaitingFor: wait.ForExit().WithExitTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if runner != nil {
		db.resources = append(db.resources, runner)
	}
	if err != nil {
		return fmt.Errorf("run liquibase: %w", err)
	}
	return checkExit(ctx, runner)
}

func checkExit(ctx context.Context, c testcontainers.Container) error {
	state, err := c.State(ctx)
	if err != nil {
		return fmt.Errorf("liquibase state: %w", err)
	}
	if state.ExitCode == 0 {
		return nil
	}
	logs, err := c.Logs(ctx)
	if err != nil {
		return fmt.Errorf("liquibase exited with %d", state.ExitCode)
	}
	defer func() { _ = logs.Close() }()
	out, _ := io.ReadAll(logs)
	return fmt.Errorf("liquibase exited with %d:\n%s", state.ExitCode, out)
}

// Stop removes every container and the network. It is safe to call twice.
func (db *Database) Stop() {
	ctx := context.Background()
	for i := len(db.resources) - 1; i >= 0; i-- {
		_ = db.resources[i].Terminate(ctx)
	}
	db.resources = nil
	if db.network != nil {
		_ = db.network.Remove(ctx)
		db.network = nil
	}
}

// databaseDir finds <repo>/database from this source file's location.
func databaseDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("locate postgrestest source")
	}
	dir := filepath.Dir(file)
	for dir != filepath.Dir(dir) {
		candidate := filepath.Join(dir, "database")
		if _, err := os.Stat(filepath.Join(candidate, "Dockerfile")); err == nil {
			return candidate, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", fmt.Errorf("database/Dockerfile not found above %s", strings.TrimSpace(filepath.Dir(file)))
}
