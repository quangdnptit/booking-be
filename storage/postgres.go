package storage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresPool is the app's PostgreSQL connection pool.
type PostgresPool = pgxpool.Pool

// PostgresConfig is built only from environment (see LoadPostgresConfigFromEnv).
type PostgresConfig struct {
	URL      string
	MaxConns int32
	MinConns int32
}

// LoadPostgresConfigFromEnv reads PostgreSQL settings from the environment only.
//
// Either set DATABASE_URL / POSTGRES_URL, or set all of:
// POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_SSLMODE.
//
// Optional: POSTGRES_MAX_CONNS, POSTGRES_MIN_CONNS.
func LoadPostgresConfigFromEnv() (PostgresConfig, error) {
	connURL := os.Getenv("DATABASE_URL")
	if connURL == "" {
		connURL = os.Getenv("POSTGRES_URL")
	}
	if connURL == "" {
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		pass := os.Getenv("POSTGRES_PASSWORD")
		db := os.Getenv("POSTGRES_DB")
		ssl := os.Getenv("POSTGRES_SSLMODE")
		switch {
		case host == "":
			return PostgresConfig{}, fmt.Errorf("postgres: POSTGRES_HOST is required (or set DATABASE_URL)")
		case port == "":
			return PostgresConfig{}, fmt.Errorf("postgres: POSTGRES_PORT is required (or set DATABASE_URL)")
		case user == "":
			return PostgresConfig{}, fmt.Errorf("postgres: POSTGRES_USER is required (or set DATABASE_URL)")
		case pass == "":
			return PostgresConfig{}, fmt.Errorf("postgres: POSTGRES_PASSWORD is required (or set DATABASE_URL)")
		case db == "":
			return PostgresConfig{}, fmt.Errorf("postgres: POSTGRES_DB is required (or set DATABASE_URL)")
		case ssl == "":
			return PostgresConfig{}, fmt.Errorf("postgres: POSTGRES_SSLMODE is required (or set DATABASE_URL)")
		}
		u := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(user, pass),
			Host:   host + ":" + port,
			Path:   "/" + db,
		}
		q := u.Query()
		q.Set("sslmode", ssl)
		u.RawQuery = q.Encode()
		connURL = u.String()
	}

	cfg := PostgresConfig{URL: connURL}
	if s := os.Getenv("POSTGRES_MAX_CONNS"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil && n > 0 {
			cfg.MaxConns = int32(n)
		}
	}
	if s := os.Getenv("POSTGRES_MIN_CONNS"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil && n >= 0 {
			cfg.MinConns = int32(n)
		}
	}
	return cfg, nil
}

// NewPostgresPool opens a pgx connection pool. Caller must Close() when done.
func NewPostgresPool(ctx context.Context, cfg PostgresConfig) (*PostgresPool, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("postgres: empty connection URL")
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres parse config: %w", err)
	}
	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute
	poolCfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	return pool, nil
}
