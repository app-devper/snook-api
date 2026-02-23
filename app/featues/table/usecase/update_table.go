package usecase

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/data/repositories"
	"snook/app/domain/request"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateTableById(tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tableId, err := primitive.ObjectIDFromHex(ctx.Param("tableId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_001, "invalid tableId")
			return
		}
		var req request.Table
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_001, err.Error())
			return
		}
		userId := ctx.GetString("UserId")
		table := entities.Table{
			Name:        req.Name,
			Type:        req.Type,
			RatePerHour: req.RatePerHour,
			Description: req.Description,
			UpdatedBy:   userId,
		}
		if err := tableEntity.UpdateTableById(tableId, table); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func UpdateTableStatus(tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tableId, err := primitive.ObjectIDFromHex(ctx.Param("tableId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_001, "invalid tableId")
			return
		}
		var req request.TableStatus
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_001, err.Error())
			return
		}
		if err := tableEntity.UpdateTableStatus(tableId, req.Status); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}
