package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"floway-backend/internal/apperr"
)

// translateNotFound maps the driver's "no rows" signal to the shared
// not-found sentinel. Every repository's single-row lookup/update/delete
// path should funnel its error through this — the driver type must never
// leak past this package.
func translateNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.ErrNotFound
	}
	return err
}

// checkDeleted turns a successful-but-zero-rows DELETE (not a pgx error —
// Postgres doesn't error when a WHERE clause matches nothing) into
// apperr.ErrNotFound, so deleting a nonexistent id doesn't silently report
// success.
func checkDeleted(tag pgconn.CommandTag, err error) error {
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
