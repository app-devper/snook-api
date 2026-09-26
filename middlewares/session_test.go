package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"snook/app/core/errcode"

	"github.com/app-devper/um-api/sessionclient"
	"github.com/gin-gonic/gin"
)

type sessionStub struct {
	userId    string
	err       error
	gotSystem string
	gotMethod string
}

func (s *sessionStub) Authorize(_ context.Context, _ string, system, method string) (string, error) {
	s.gotSystem, s.gotMethod = system, method
	return s.userId, s.err
}

func runRequireSession(stub *sessionStub, method string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(method, "/", nil)
	ctx.Set("SessionId", "session-1")
	ctx.Set("System", "SNOOK")
	RequireSession(stub)(ctx)
	return ctx, w
}

func TestRequireSessionSetsUserIdFromLiveSession(t *testing.T) {
	stub := &sessionStub{userId: "user-1"}
	ctx, w := runRequireSession(stub, http.MethodGet)
	if ctx.IsAborted() || ctx.GetString("UserId") != "user-1" {
		t.Fatalf("expected UserId user-1, got %q (status %d)", ctx.GetString("UserId"), w.Code)
	}
	if stub.gotSystem != "SNOOK" || stub.gotMethod != http.MethodGet {
		t.Fatalf("expected token system and request method, got %q %q", stub.gotSystem, stub.gotMethod)
	}
}

func TestRequireSessionRejectsRevokedSession(t *testing.T) {
	_, w := runRequireSession(&sessionStub{err: sessionclient.ErrSessionRejected}, http.MethodGet)
	if w.Code != http.StatusUnauthorized || !strings.Contains(w.Body.String(), errcode.AU_UNAUTHORIZED_005) {
		t.Fatalf("expected 401 %s, got %d %s", errcode.AU_UNAUTHORIZED_005, w.Code, w.Body.String())
	}
}

func TestRequireSessionReturns503WhenSessionStoreUnavailable(t *testing.T) {
	_, w := runRequireSession(&sessionStub{err: sessionclient.ErrUnavailable}, http.MethodPost)
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(w.Body.String(), errcode.AU_UNAVAILABLE_001) {
		t.Fatalf("expected 503 %s, got %d %s", errcode.AU_UNAVAILABLE_001, w.Code, w.Body.String())
	}
}
