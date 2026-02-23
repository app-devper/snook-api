package setting

import (
	"net/http"
	"snook/app/core/constant"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/domain"
	"snook/app/domain/request"
	"snook/middlewares"

	"github.com/gin-gonic/gin"
)

func ApplySettingAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("settings")

	r.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		s, err := repository.Setting.GetSetting()
		if err != nil {
			ctx.JSON(http.StatusOK, entities.Setting{})
			return
		}
		ctx.JSON(http.StatusOK, s)
	})

	r.PUT("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			var req request.Setting
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.SE_BAD_REQUEST_001, err.Error())
				return
			}
			s := entities.Setting{
				CompanyName: req.CompanyName, CompanyAddress: req.CompanyAddress,
				CompanyPhone: req.CompanyPhone, CompanyTaxId: req.CompanyTaxId,
				ReceiptFooter: req.ReceiptFooter, PromptPayId: req.PromptPayId,
				UpdatedBy: ctx.GetString("UserId"),
			}
			if err := repository.Setting.UpsertSetting(s); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.SE_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})
}
