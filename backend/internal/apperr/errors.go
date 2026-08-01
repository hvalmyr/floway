// Package apperr holds the sentinel errors shared by every service and
// repository. Repositories translate driver-specific errors (pgx.ErrNoRows,
// unique-violation codes) into these at their own boundary — the service and
// handler layers never see a driver type. See internal/service for the
// service.ErrValidation/ErrNotFound aliases kept for call-site compatibility.
package apperr

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)
