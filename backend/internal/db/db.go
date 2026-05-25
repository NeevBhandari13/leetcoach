package db

import (
	"database/sql"
	"embed"
	"encoding/json"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

//go:embed seed/problems.json
var problemsSeedJSON []byte

// dsn is the data source name
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	// check if error with initialisation
	if err != nil {
		return nil, err
	}
	// check if we can ping the db
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(sqlDB *sql.DB) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

type problem struct {
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	Difficulty  string          `json:"difficulty"`
	Description string          `json:"description"`
	Examples    json.RawMessage `json:"examples"`
	Constraints string          `json:"constraints"`
	Topics      []string        `json:"topics"`
}

// Seed inserts problems from the embedded JSON file. It is idempotent —
// rows that already exist (matched by slug) are skipped.
func Seed(sqlDB *sql.DB) error {
	var problems []problem
	if err := json.Unmarshal(problemsSeedJSON, &problems); err != nil {
		return err
	}

	for _, p := range problems {
		_, err := sqlDB.Exec(`
			INSERT INTO problems (slug, title, difficulty, description, examples, constraints, topics)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (slug) DO NOTHING`,
			p.Slug, p.Title, p.Difficulty, p.Description, p.Examples, p.Constraints, pq.Array(p.Topics),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
