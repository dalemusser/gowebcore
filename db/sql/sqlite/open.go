package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Open returns a *sqlx.DB for an in-memory or file path DSN.
func Open(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dsn) // e.g. "file:test.db?_busy_timeout=5000"
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
