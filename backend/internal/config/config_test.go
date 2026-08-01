package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/config"
)

// setRequiredEnv sets every OTHER required var so a test can isolate the
// one it's actually exercising.
func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "local")
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("GARAGE_ENDPOINT", "http://garage:3900")
	t.Setenv("GARAGE_ACCESS_KEY", "key")
	t.Setenv("GARAGE_SECRET_KEY", "secret")
	t.Setenv("GARAGE_BUCKET", "bucket")
}

func TestLoad_ValidConfig(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "local", cfg.Env)
}

// A missing APP_ENV must never silently fall back to "local" — that default
// also drives SecureCookies (cfg.Env != "local"), so a production deploy
// that forgets to set it would silently ship session cookies without the
// Secure flag (architecture review finding #9).
func TestLoad_RequiresAppEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("APP_ENV", "")

	_, err := config.Load()

	require.Error(t, err)
}

func TestLoad_RejectsUnknownAppEnv(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("APP_ENV", "prod") // common typo for "production"

	_, err := config.Load()

	require.Error(t, err)
}

func TestLoad_RequiresDatabaseURL(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DATABASE_URL", "")

	_, err := config.Load()

	require.Error(t, err)
}
