// Package postgrestest starts a throwaway PostgreSQL for integration tests,
// with every Liquibase migration applied by the same runner image used in
// compose and CI (database/Dockerfile). It requires Docker and the runner
// image built beforehand (task backend:test does it): building it inside
// every test process ran the same build and driver download in parallel.
package postgrestest

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/containerd/errdefs"
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
	// MigrationsImage is the Liquibase runner image the tests apply the
	// schema with. Build it from the repository root with:
	//
	//	docker build -t openboard-migrations:test database
	MigrationsImage = "openboard-migrations:test"
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
	if err := requireImage(ctx, MigrationsImage); err != nil {
		return err
	}
	runner, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:    MigrationsImage,
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

// requireImage fails fast when the image was not built locally; otherwise
// testcontainers would try, and retry, to pull it from a registry.
func requireImage(ctx context.Context, name string) error {
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer func() { _ = cli.Close() }()
	_, err = cli.ImageInspect(ctx, name)
	if errdefs.IsNotFound(err) {
		return fmt.Errorf("image %s not found: run `task backend:test`, or `docker build -t %s database` from the repository root", name, name)
	}
	if err != nil {
		return fmt.Errorf("inspect image %s: %w", name, err)
	}
	return nil
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
