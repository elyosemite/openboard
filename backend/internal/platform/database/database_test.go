package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/elyosemite/openboard/backend/internal/testsupport/postgrestest"
)

var dbURL string

func TestMain(m *testing.M) { os.Exit(runWithDatabase(m)) }

func runWithDatabase(m *testing.M) int {
	db, err := postgrestest.Start(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer db.Stop()
	dbURL = db.URL
	return m.Run()
}

func TestOpenAndCheck(t *testing.T) {
	pool, err := Open(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer pool.Close()

	c := Checker{Pool: pool}
	if c.Name() != "postgres" {
		t.Errorf("Name = %q", c.Name())
	}
	if err := c.Check(context.Background()); err != nil {
		t.Errorf("Check: %v", err)
	}

	var users int
	if err := pool.QueryRow(context.Background(), "SELECT count(*) FROM users").Scan(&users); err != nil {
		t.Errorf("migrations should have created the users table: %v", err)
	}
}

func TestOpenFailures(t *testing.T) {
	cases := map[string]string{
		"malformed url": "postgres://%zz",
		"unreachable":   "postgres://openboard:openboard@127.0.0.1:1/openboard?sslmode=disable&connect_timeout=2",
	}
	for name, url := range cases {
		if _, err := Open(context.Background(), url); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}
