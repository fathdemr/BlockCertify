package helper

import (
	"BlockCertify/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name string, value string, expiration time.Time) {
	maxAge := int(time.Until(expiration).Seconds())
	if maxAge < 0 {
		maxAge = -1
	}
	http.SetCookie(c.Writer, buildCookie(name, value, maxAge))
}

func ClearCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, buildCookie(name, "", -1))
}

// buildCookie hardens auth cookies for browser use:
//   - HttpOnly: not readable from JavaScript (XSS can't steal the token)
//   - Secure (live only): never sent over plain HTTP; disabled locally so
//     http://localhost development keeps working
//   - SameSite=Lax: browsers won't attach the cookie to cross-site POSTs,
//     which blocks CSRF against the cookie-authenticated endpoints
func buildCookie(name string, value string, maxAge int) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.Params.GetString("environment") == "live",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
	return cookie
}
