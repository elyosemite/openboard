package postgrestest

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestStartAppliesEveryMigration(t *testing.T) {
	ctx := context.Background()
	db, err := Start(ctx)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer db.Stop()

	conn, err := pgx.Connect(ctx, db.URL)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	var applied int
	if err := conn.QueryRow(ctx, `SELECT count(*) FROM databasechangelog`).Scan(&applied); err != nil {
		t.Fatalf("Liquibase should have recorded its changesets: %v", err)
	}
	if applied == 0 {
		t.Error("no changeset was applied")
	}
	if _, err := conn.Exec(ctx, `SELECT 1 FROM users LIMIT 1`); err != nil {
		t.Errorf("users table missing: %v", err)
	}
}

func TestStopIsIdempotent(t *testing.T) {
	db := &Database{}
	db.Stop()
	db.Stop()
	if db.resources != nil || db.network != nil {
		t.Error("Stop should clear every resource")
	}
}

func TestRequireImageExplainsHowToBuildAMissingImage(t *testing.T) {
	err := requireImage(context.Background(), "openboard-migrations:does-not-exist")
	if err == nil {
		t.Fatal("a missing image should be an error")
	}
	if !strings.Contains(err.Error(), "task backend:test") {
		t.Errorf("error should say how to build the image, got: %v", err)
	}
}
