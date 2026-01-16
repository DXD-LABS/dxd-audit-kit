package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
)

//go:embed *.sql
var migrationsFS embed.FS

// RunMigrations thực hiện chạy các file migration sql từ thư mục migrations
func RunMigrations(db *sql.DB) error {
	// Tạo bảng schema_migrations nếu chưa có
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	entries, err := migrationsFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name()[len(entry.Name())-4:] == ".sql" {
			filenames = append(filenames, entry.Name())
		}
	}
	sort.Strings(filenames)

	for _, filename := range filenames {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", filename, err)
		}

		if exists {
			continue
		}

		fmt.Printf("Running migration: %s\n", filename)
		content, err := migrationsFS.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}
	}

	return nil
}
