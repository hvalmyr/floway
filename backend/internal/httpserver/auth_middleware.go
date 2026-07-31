package httpserver

import (
	"context"
	"errors"
	"net/http"

	"floway-backend/internal/auth"
)

const sessionCookieName = "floway_admin_session"

var errUnauthorized = errors.New("unauthorized")

type adminContextKey struct{}

type adminContext struct {
	UserID int64
	Login  string
}

func requireAdminMiddleware(tokens *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil {
				writeError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			claims, err := tokens.Parse(cookie.Value)
			if err != nil {
				writeError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), adminContextKey{}, adminContext{
				UserID: claims.UserID,
				Login:  claims.Login,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
