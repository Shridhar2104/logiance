// internal/database/migrate/migrate.go
package migrate

import (
    "database/sql"
    "errors"
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create migration driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", migrationsPath),
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("could not create migrate instance: %w", err)
    }

    err = m.Up()
    if err != nil && !errors.Is(err, migrate.ErrNoChange) {
        return fmt.Errorf("could not run migrations: %w", err)
    }

    if errors.Is(err, migrate.ErrNoChange) {
        log.Println("No migrations to run")
        return nil
    }

    log.Println("Migrations completed successfully")
    return nil
}