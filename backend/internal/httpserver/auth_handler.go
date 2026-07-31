package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"floway-backend/internal/auth"
	"floway-backend/internal/service"
)

type authHandler struct {
	svc           *service.AdminUserService
	tokens        *auth.TokenManager
	secureCookies bool
	admin         func(http.Handler) http.Handler
}

func newAuthHandler(svc *service.AdminUserService, tokens *auth.TokenManager, secureCookies bool, admin func(http.Handler) http.Handler) *authHandler {
	return &authHandler{svc: svc, tokens: tokens, secureCookies: secureCookies, admin: admin}
}

func (h *authHandler) routes(r chi.Router) {
	r.Post("/login", h.login)
	r.Post("/logout", h.logout)
	r.With(h.admin).Get("/me", h.me)
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.svc.Authenticate(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	token, expiresAt, err := h.tokens.Issue(user.ID, user.Login)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	h.setSessionCookie(w, token, expiresAt)
	writeJSON(w, http.StatusOK, map[string]string{"login": user.Login})
}

func (h *authHandler) logout(w http.ResponseWriter, r *http.Request) {
	h.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *authHandler) me(w http.ResponseWriter, r *http.Request) {
	adm, ok := r.Context().Value(adminContextKey{}).(adminContext)
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": adm.UserID, "login": adm.Login})
}

func (h *authHandler) setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *authHandler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}
