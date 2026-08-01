package service

import (
	"errors"

	"floway-backend/internal/apperr"
)

// Aliased from internal/apperr rather than redefined here: repositories
// return apperr.ErrNotFound/apperr.ErrValidation directly (they're the same
// error value), so every existing errors.Is(err, service.ErrNotFound) check
// across handlers and tests keeps working unchanged.
var (
	ErrValidation = apperr.ErrValidation
	ErrNotFound   = apperr.ErrNotFound
	ErrConflict   = apperr.ErrConflict
)

// ErrInvalidCredentials is returned by AdminUserService.Authenticate when
// the login does not exist or the password does not match — intentionally
// non-disclosing about which of the two it is.
var ErrInvalidCredentials = errors.New("invalid credentials")
