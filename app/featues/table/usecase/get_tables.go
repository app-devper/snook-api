package usecase

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/repositories"

	"github.com/gin-gonic/gin"
)

func GetTables(tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tables, err := tableEntity.GetTables()
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.TB_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, tables)
	}
}
