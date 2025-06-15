package migrate

import (
	"context"
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Run applies all pending migrations.
func Run(ctx context.Context, db *sqlx.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect(db.DriverName()); err != nil {
		return err
	}
	return goose.UpContext(ctx, db.DB, "migrations") // note db.DB
}
