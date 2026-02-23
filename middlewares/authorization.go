package middlewares

import (
	"net/http"
	"snook/app/core/errcode"

	"github.com/gin-gonic/gin"
)

func RequireAuthorization(auths ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("Role")
		if role == "" {
			invalidRequest(ctx)
			return
		}
		isAccessible := false
		for _, auth := range auths {
			if role == auth {
				isAccessible = true
				break
			}
		}
		if !isAccessible {
			notPermission(ctx)
			return
		}
		ctx.Next()
	}
}

func invalidRequest(ctx *gin.Context) {
	errcode.Abort(ctx, http.StatusForbidden, errcode.SY_FORBIDDEN_001, "Invalid request, restricted endpoint")
}

func notPermission(ctx *gin.Context) {
	errcode.Abort(ctx, http.StatusForbidden, errcode.SY_FORBIDDEN_002, "Don't have permission")
}
