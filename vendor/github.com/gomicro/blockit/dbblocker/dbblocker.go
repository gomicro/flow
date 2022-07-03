package dbblocker

import (
	"context"
	"database/sql"
	"time"

	"github.com/gomicro/blockit/cbblocker"
)

// New takes a SQL database object and returns a newly instantiated Blocker
func New(db *sql.DB) *cbblocker.Blocker {
	return cbblocker.New(db.Ping, 1*time.Second)
}

// NewWithContext taxes a context and SQL dtabase object and returns a newly
// instantiated Blocker with the context passed along to the database object.
func NewWithContext(ctx context.Context, db *sql.DB) *cbblocker.Blocker {
	return cbblocker.New(func() error {
		return db.PingContext(ctx)
	}, 1*time.Second)
}
