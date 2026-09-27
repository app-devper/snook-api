package middlewares

import (
	"errors"
	"os"
	"snook/app/core/errcode"

	"github.com/app-devper/um-api/sessionclient"
	"github.com/app-devper/um-api/sessionclient/ginauth"
	"github.com/gin-gonic/gin"
)

// NewAuth verifies UM access tokens for snook: SYSTEM and CLIENT_ID pin the
// token, and the session is confirmed in UM's Redis at redisHost
// (um-api ADR-0005). It fails when any of them is missing.
func NewAuth(redisHost string) (*ginauth.Auth, error) {
	store, err := sessionclient.RedisStoreFor(redisHost)
	if err != nil {
		return nil, err
	}
	return NewAuthWithStore(store)
}

// NewAuthWithStore is NewAuth with UM's session store supplied, for tests.
func NewAuthWithStore(store sessionclient.Store) (*ginauth.Auth, error) {
	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("missing required env: CLIENT_ID")
	}
	verifier, err := sessionclient.NewVerifier(sessionclient.Config{
		SecretKey: os.Getenv("SECRET_KEY"),
		System:    os.Getenv("SYSTEM"),
		ClientID:  clientId,
		Store:     store,
	})
	if err != nil {
		return nil, err
	}
	return ginauth.New(verifier, func(ctx *gin.Context, e *sessionclient.Error) {
		errcode.Abort(ctx, e.Status, e.Code, e.Message)
	}), nil
}

// RequireSession admits a caller with a live UM session. Every snook route
// uses the default outage policy: while UM is unreachable a read may continue
// with the session last confirmed for its token, and writes wait.
func RequireSession(auth *ginauth.Auth) gin.HandlerFunc {
	return auth.Require(sessionclient.ReadOnlyWithLastGood)
}
