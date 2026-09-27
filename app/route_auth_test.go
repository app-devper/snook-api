package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"snook/app/domain"
	"snook/middlewares"

	"github.com/app-devper/um-api/sessionclient"
	"github.com/gin-gonic/gin"
)

// unusedStore stands in for UM's session store; anonymous requests never
// reach it.
type unusedStore struct{}

func (unusedStore) Session(context.Context, string) (sessionclient.Session, error) {
	return sessionclient.Session{}, sessionclient.ErrUnavailable
}

// Every business route must refuse a request without a UM token before it
// reaches a repository, so the nil repositories here are never used.
func TestEveryBusinessRouteRejectsAnAnonymousRequest(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("CLIENT_ID", "000")
	t.Setenv("SYSTEM", "SNOOK")
	auth, err := middlewares.NewAuthWithStore(unusedStore{})
	if err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	applyFeatureAPIs(r.Group("/api/snook/v1"), &domain.Repository{Auth: auth})

	checked := 0
	for _, route := range r.Routes() {
		checked++
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(route.Method, strings.ReplaceAll(route.Path, ":", "x"), nil))
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s: expected 401 without a token, got %d", route.Method, route.Path, w.Code)
		}
	}
	if checked == 0 {
		t.Fatal("no routes registered")
	}
	t.Logf("checked %d business routes", checked)
}

func TestNewAuthRequiresClientId(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("SYSTEM", "SNOOK")
	t.Setenv("CLIENT_ID", "")
	if _, err := middlewares.NewAuthWithStore(unusedStore{}); err == nil {
		t.Fatal("expected an error without CLIENT_ID")
	}
}
