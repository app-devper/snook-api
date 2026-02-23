package usecase

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/data/repositories"
	"snook/app/domain/request"

	"github.com/gin-gonic/gin"
)

func CreateTable(tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req request.Table
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_001, err.Error())
			return
		}
		userId := ctx.GetString("UserId")
		table := entities.Table{
			Name:        req.Name,
			Type:        req.Type,
			Status:      "AVAILABLE",
			RatePerHour: req.RatePerHour,
			Description: req.Description,
			CreatedBy:   userId,
		}
		result, err := tableEntity.CreateTable(table)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TB_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, result)
	}
}
