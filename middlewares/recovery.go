package middlewares

import (
	"net/http"
	"snook/app/core/errcode"

	"github.com/gin-gonic/gin"
)

func NewRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		errcode.Abort(ctx, http.StatusInternalServerError, errcode.SY_INTERNAL_001, "internal server error")
	})
}
