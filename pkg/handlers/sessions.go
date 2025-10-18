package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"catecard/pkg/domain/entities"
)

var (
	// cookie name
	sessionCookieName = "session_token"
	// db handle will be initialized from web.go
	sessionsDB *sql.DB
)

// InitSessionStore inicializa el DB handle que usa el store de sesiones.
func InitSessionStore(db *sql.DB) {
	sessionsDB = db
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateSession guarda la sesión en la tabla sessions y escribe la cookie.
func CreateSession(w http.ResponseWriter, user *entities.User) (string, error) {
	token, err := generateToken(32)
	if err != nil {
		return "", err
	}

	// si no hay DB, fallback a cookie sin persistencia (comportamiento previo)
	// If sessionsDB is nil (fallback), still set cookie; allow disabling Secure for local/dev via env
	if sessionsDB == nil {
		cookie := &http.Cookie{
			Name:     sessionCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   cookieSecure(),
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(24 * time.Hour),
		}
		if d := cookieDomain(); d != "" {
			cookie.Domain = d
		}
		http.SetCookie(w, cookie)
		return token, nil
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	expires := time.Now().Add(24 * time.Hour)
	_, err = sessionsDB.Exec(`INSERT INTO sessions(token, user_json, expires_at) VALUES(?, ?, ?)`, token, string(userJSON), expires)
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   cookieSecure(),
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	}
	if d := cookieDomain(); d != "" {
		cookie.Domain = d
	}
	http.SetCookie(w, cookie)
	return token, nil
}

// DeleteSession borra la sesión en DB (si existe) y expira la cookie en el cliente.
func DeleteSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(sessionCookieName)
	if err == nil && sessionsDB != nil {
		token := c.Value
		sessionsDB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	}

	expired := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   cookieSecure(),
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	if d := cookieDomain(); d != "" {
		expired.Domain = d
	}
	http.SetCookie(w, expired)
}

// cookieSecure returns whether session cookies should use Secure flag.
// For local HTTP development you can disable it by setting DEV_DISABLE_SECURE_COOKIE=1
func cookieSecure() bool {
	return os.Getenv("DEV_DISABLE_SECURE_COOKIE") != "1"
}

// cookieDomain returns the domain attribute for the session cookie if set via env
func cookieDomain() string {
	return os.Getenv("COOKIE_DOMAIN")
}

// GetUserFromRequest devuelve el usuario asociado a la cookie (o nil).
func GetUserFromRequest(r *http.Request) *entities.User {
	c, err := r.Cookie(sessionCookieName)
	if err != nil || sessionsDB == nil {
		return nil
	}
	token := c.Value

	row := sessionsDB.QueryRow(`SELECT user_json, expires_at FROM sessions WHERE token = ?`, token)
	var userJSON string
	var expiresAt time.Time
	if err := row.Scan(&userJSON, &expiresAt); err != nil {
		return nil
	}
	if time.Now().After(expiresAt) {
		// session expired
		sessionsDB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
		return nil
	}
	var u entities.User
	if err := json.Unmarshal([]byte(userJSON), &u); err != nil {
		return nil
	}
	return &u
}
