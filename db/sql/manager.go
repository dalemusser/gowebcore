package sql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Manager struct {
	conns map[string]*sqlx.DB
}

func New() *Manager { return &Manager{conns: make(map[string]*sqlx.DB)} }

// Add stores a *sqlx.DB under a key (e.g., "primary").
func (m *Manager) Add(name string, db *sqlx.DB) { m.conns[name] = db }

// DB returns the connection by name.
func (m *Manager) DB(name string) *sqlx.DB { return m.conns[name] }

// CloseAll closes every connection.
func (m *Manager) CloseAll(ctx context.Context) error {
	for _, db := range m.conns {
		_ = db.Close()
	}
	return nil
}

/* ---------- convenience helpers (optional) ---------------------------- */

func Exec(ctx context.Context, db *sqlx.DB, query string, args ...any) error {
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func Ping(db *sqlx.DB) error { return db.Ping() }

func SetMax(db *sqlx.DB) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
}
