package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const refreshTokenCookieName = "fortune_refresh_token"

func setRefreshTokenCookie(c *gin.Context, token string, expiresAt int64, remember bool) {
	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   requestIsSecure(c.Request),
		SameSite: http.SameSiteLaxMode,
	}
	if remember {
		cookie.Expires = time.Unix(expiresAt, 0)
		cookie.MaxAge = int(time.Until(cookie.Expires).Seconds())
		if cookie.MaxAge < 1 {
			cookie.MaxAge = 1
		}
	}
	http.SetCookie(c.Writer, cookie)
}

func clearRefreshTokenCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshTokenCookieName,
		Path:     "/api/auth",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   requestIsSecure(c.Request),
		SameSite: http.SameSiteLaxMode,
	})
}

func requestIsSecure(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}
