package usecase

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTableById(tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tableId, err := primitive.ObjectIDFromHex(ctx.Param("tableId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_001, "invalid tableId")
			return
		}
		table, err := tableEntity.GetTableById(tableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_002, "table not found")
			return
		}
		if table.Status == "IN_USE" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_002, "cannot delete table while in use")
			return
		}
		if err := tableEntity.DeleteTableById(tableId); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}
