package helper

import (
	"BlockCertify/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setCookieAndParse(t *testing.T, expiration time.Time) *http.Cookie {
	t.Helper()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SetCookie(c, "jwt", "tokenvalue", expiration)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	return cookies[0]
}

func TestSetCookieMaxAgeMatchesExpiration(t *testing.T) {
	cookie := setCookieAndParse(t, time.Now().Add(time.Hour))

	// Allow a couple of seconds of slack for test execution time.
	if cookie.MaxAge < 3595 || cookie.MaxAge > 3600 {
		t.Errorf("MaxAge = %d, want ~3600", cookie.MaxAge)
	}
}

func TestSetCookieSecurityFlags(t *testing.T) {
	cookie := setCookieAndParse(t, time.Now().Add(time.Hour))

	if !cookie.HttpOnly {
		t.Error("cookie is not HttpOnly")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("SameSite = %v, want Lax", cookie.SameSite)
	}
}

func TestSetCookieSecureInLiveEnvironment(t *testing.T) {
	original := config.Params.GetString("environment")
	config.Params.Set("environment", "live")
	t.Cleanup(func() { config.Params.Set("environment", original) })

	cookie := setCookieAndParse(t, time.Now().Add(time.Hour))
	if !cookie.Secure {
		t.Error("cookie is not Secure in live environment")
	}
}

func TestClearCookieExpiresImmediately(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ClearCookie(c, "jwt")

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].MaxAge >= 0 {
		t.Errorf("MaxAge = %d, want negative (delete)", cookies[0].MaxAge)
	}
	if cookies[0].Value != "" {
		t.Errorf("Value = %q, want empty", cookies[0].Value)
	}
}
