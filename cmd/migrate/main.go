package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	"todo/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fatal("failed to load configuration: " + err.Error())
	}

	migrationsDir := flag.String("dir", "", "migrations directory (default: auto-detect)")
	flag.Parse()

	dbURL := resolveDBURL(cfg)
	if dbURL == "" {
		fatal("MIGRATIONS_DATABASE_URL, DATABASE_URL, or valid DB config is required")
	}

	dir := *migrationsDir
	if dir == "" {
		dir = findMigrationsDir()
	}

	args := flag.Args()
	if len(args) == 0 {
		fatal("usage: migrate <up|down|status|force|create> [args]")
	}

	switch args[0] {
	case "up":
		runUp(dbURL, dir)
	case "down":
		n := 1
		if len(args) > 1 {
			var err error
			n, err = strconv.Atoi(args[1])
			if err != nil {
				fatal("down: invalid step count: " + args[1])
			}
		}
		runDown(dbURL, dir, n)
	case "status":
		runStatus(dbURL, dir)
	case "force":
		if len(args) < 2 {
			fatal("usage: migrate force <version>")
		}
		v, err := strconv.Atoi(args[1])
		if err != nil {
			fatal("force: invalid version: " + args[1])
		}
		runForce(dbURL, dir, v)
	case "create":
		if len(args) < 2 {
			fatal("usage: migrate create <name>")
		}
		runCreate(dir, args[1])
	default:
		fatal("unknown command: " + args[0])
	}
}

func newMigrate(dbURL, dir string) *migrate.Migrate {
	m, err := migrate.New("file://"+dir, dbURL)
	if err != nil {
		fatal("failed to init migrate: " + err.Error())
	}
	return m
}

func runUp(dbURL, dir string) {
	if err := ensureDatabaseExists(dbURL); err != nil {
		fatal("ensure database: " + err.Error())
	}
	m := newMigrate(dbURL, dir)
	defer func() { _, _ = m.Close() }()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fatal("up: " + err.Error())
	}
	fmt.Println("migrations applied successfully")
}

func ensureDatabaseExists(dbURL string) error {
	u, err := url.Parse(dbURL)
	if err != nil {
		return fmt.Errorf("invalid database URL: %w", err)
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" || dbName == "postgres" {
		return nil
	}

	// Try connecting to default admin database "postgres", then fallback to "template1"
	adminURL := *u
	adminURL.Path = "/postgres"

	adminDB, err := sql.Open("postgres", adminURL.String())
	if err != nil {
		return fmt.Errorf("failed to open admin database connection: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		// Fallback to template1 if postgres db is not accessible
		adminURL.Path = "/template1"
		altDB, altErr := sql.Open("postgres", adminURL.String())
		if altErr != nil || altDB.Ping() != nil {
			if altDB != nil {
				_ = altDB.Close()
			}
			return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
		}
		_ = adminDB.Close()
		adminDB = altDB
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	if err := adminDB.QueryRow(query, dbName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		createStmt := fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName))
		if _, err := adminDB.Exec(createStmt); err != nil {
			return fmt.Errorf("failed to create database %q: %w", dbName, err)
		}
		fmt.Printf("database %q created successfully\n", dbName)
	}

	return nil
}

func runDown(dbURL, dir string, steps int) {
	m := newMigrate(dbURL, dir)
	defer func() { _, _ = m.Close() }()
	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fatal("down: " + err.Error())
	}
	fmt.Printf("rolled back %d migration(s)\n", steps)
}

func runForce(dbURL, dir string, version int) {
	m := newMigrate(dbURL, dir)
	defer func() { _, _ = m.Close() }()
	if err := m.Force(version); err != nil {
		fatal("force: " + err.Error())
	}
	fmt.Printf("forced version to %d\n", version)
}

func runStatus(dbURL, dir string) {
	m := newMigrate(dbURL, dir)
	defer func() { _, _ = m.Close() }()
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("no migrations applied yet")
			return
		}
		fatal("status: " + err.Error())
	}
	dirtyStr := ""
	if dirty {
		dirtyStr = " (DIRTY)"
	}
	fmt.Printf("current version: %d%s\n", version, dirtyStr)
}

func runCreate(dir, name string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fatal("create: read dir: " + err.Error())
	}
	next := 1
	for _, e := range entries {
		if len(e.Name()) < 3 {
			continue
		}
		if n, err := strconv.Atoi(e.Name()[:3]); err == nil && n >= next {
			next = n + 1
		}
	}
	prefix := fmt.Sprintf("%03d_%s", next, name)
	for _, suffix := range []string{".up.sql", ".down.sql"} {
		path := filepath.Join(dir, prefix+suffix)
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			fatal("create: " + err.Error())
		}
		fmt.Println("created:", path)
	}
}

func findMigrationsDir() string {
	candidates := []string{
		"migrations",
		"../migrations",
		"../../migrations",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	fatal("could not find migrations directory — use -dir flag")
	return ""
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "migrate:", msg)
	os.Exit(1)
}

func resolveDBURL(cfg *config.Config) string {
	if v := os.Getenv("MIGRATIONS_DATABASE_URL"); v != "" {
		return v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	if cfg != nil {
		return cfg.Database.DSN()
	}
	return ""
}
