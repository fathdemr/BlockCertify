package middleware

import (
	"BlockCertify/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type fakeTokenHelper struct {
	claims jwt.MapClaims
	err    error
}

func (f *fakeTokenHelper) Sign(claims jwt.MapClaims) (string, error) { return "token", nil }
func (f *fakeTokenHelper) Verify(tokenString string) (jwt.MapClaims, error) {
	return f.claims, f.err
}

type fakeUserRepo struct {
	user *models.User
	err  error
}

func (f *fakeUserRepo) FindByEmail(email string) (*models.User, error) { return f.user, f.err }
func (f *fakeUserRepo) Exists(email string) (bool, error)             { return f.user != nil, nil }
func (f *fakeUserRepo) Create(user *models.User) error                { return nil }
func (f *fakeUserRepo) CreateAdmin(admin *models.Admin) error         { return nil }
func (f *fakeUserRepo) CreateTransaction() *gorm.DB                   { return nil }

func setupRouter(auth AuthMiddleware, role models.UserRole) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/protected", auth.Authorize(), auth.RequireRole(role), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestAuthorizeWithoutTokenReturns401(t *testing.T) {
	auth := NewAuthMiddleware(&fakeTokenHelper{}, &fakeUserRepo{})
	r := setupRouter(auth, models.RoleAdmin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestAuthorizeWithUnknownUserReturns401(t *testing.T) {
	auth := NewAuthMiddleware(
		&fakeTokenHelper{claims: jwt.MapClaims{"email": "ghost@example.com"}},
		&fakeUserRepo{err: gorm.ErrRecordNotFound},
	)
	r := setupRouter(auth, models.RoleAdmin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestRequireRoleRejectsWrongRole(t *testing.T) {
	auth := NewAuthMiddleware(
		&fakeTokenHelper{claims: jwt.MapClaims{"email": "student@example.com"}},
		&fakeUserRepo{user: &models.User{Email: "student@example.com", Role: models.RoleStudent}},
	)
	r := setupRouter(auth, models.RoleAdmin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Code)
	}
}

func TestRequireRoleAllowsCorrectRole(t *testing.T) {
	auth := NewAuthMiddleware(
		&fakeTokenHelper{claims: jwt.MapClaims{"email": "admin@example.com"}},
		&fakeUserRepo{user: &models.User{Email: "admin@example.com", Role: models.RoleAdmin}},
	)
	r := setupRouter(auth, models.RoleAdmin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("Authorization", "Bearer sometoken")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestAuthorizeReadsJWTFromCookie(t *testing.T) {
	auth := NewAuthMiddleware(
		&fakeTokenHelper{claims: jwt.MapClaims{"email": "admin@example.com"}},
		&fakeUserRepo{user: &models.User{Email: "admin@example.com", Role: models.RoleAdmin}},
	)
	r := setupRouter(auth, models.RoleAdmin)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: "sometoken"})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}
