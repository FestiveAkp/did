package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func dbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".did")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "did.db"), nil
}

func openDB() (*sql.DB, error) {
	path, err := dbPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// migrate applies migrations/*.sql in filename order, tracked via
// PRAGMA user_version. Add a new numbered file (e.g. 0002_create_x.sql)
// for future schema changes rather than editing existing ones.
func migrate(db *sql.DB) error {
	entries, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)

	var version int
	if err := db.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}

	for ; version < len(entries); version++ {
		sqlBytes, err := migrationsFS.ReadFile(entries[version])
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entries[version], err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", entries[version], err)
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", entries[version], err)
		}
		if _, err := tx.Exec(fmt.Sprintf(`PRAGMA user_version = %d`, version+1)); err != nil {
			tx.Rollback()
			return fmt.Errorf("set schema version %d: %w", version+1, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", entries[version], err)
		}
	}

	return nil
}
