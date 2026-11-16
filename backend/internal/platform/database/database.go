// Package database opens the PostgreSQL connection pool and reports its health.
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Open creates a connection pool for url and verifies it can reach the server.
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("configure pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("reach database: %w", err)
	}
	return pool, nil
}

// Checker is a readiness check that pings the database.
type Checker struct {
	Pool *pgxpool.Pool
}

// Name identifies the check in logs.
func (Checker) Name() string { return "postgres" }

// Check pings the database.
func (c Checker) Check(ctx context.Context) error { return c.Pool.Ping(ctx) }
