package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"floway-backend/internal/auth"
)

func TestTokenManager_IssueAndParse(t *testing.T) {
	tm := auth.NewTokenManager("test-secret", time.Hour)

	token, expiresAt, err := tm.Issue(42, "admin", 3)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(time.Hour), expiresAt, time.Second)

	claims, err := tm.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, int64(42), claims.UserID)
	assert.Equal(t, "admin", claims.Login)
	assert.Equal(t, 3, claims.Version)
}

func TestTokenManager_Parse_RejectsExpiredToken(t *testing.T) {
	tm := auth.NewTokenManager("test-secret", -time.Hour)

	token, _, err := tm.Issue(1, "admin", 0)
	require.NoError(t, err)

	_, err = tm.Parse(token)
	require.Error(t, err)
	assert.ErrorIs(t, err, auth.ErrExpiredToken)
}

func TestTokenManager_Parse_RejectsWrongSecret(t *testing.T) {
	issuer := auth.NewTokenManager("secret-a", time.Hour)
	verifier := auth.NewTokenManager("secret-b", time.Hour)

	token, _, err := issuer.Issue(1, "admin", 0)
	require.NoError(t, err)

	_, err = verifier.Parse(token)
	require.Error(t, err)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestTokenManager_Parse_RejectsGarbage(t *testing.T) {
	tm := auth.NewTokenManager("test-secret", time.Hour)

	_, err := tm.Parse("not-a-jwt")
	require.Error(t, err)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}
