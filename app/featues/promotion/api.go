package promotion

import (
	"net/http"
	"snook/app/core/constant"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/domain"
	"snook/app/domain/request"
	"snook/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyPromotionAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("promotions")

	r.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		promos, err := repository.Promotion.GetPromotions()
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.PM_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, promos)
	})

	r.GET("/active", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		tableType := ctx.Query("tableType")
		promos, err := repository.Promotion.GetActivePromotions(tableType)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.PM_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, promos)
	})

	r.GET("/:promotionId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("promotionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_001, "invalid promotionId")
			return
		}
		promo, err := repository.Promotion.GetPromotionById(id)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, promo)
	})

	r.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			var req request.Promotion
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_001, err.Error())
				return
			}
			startDate, _ := time.Parse("2006-01-02", req.StartDate)
			endDate, _ := time.Parse("2006-01-02", req.EndDate)
			status := req.Status
			if status == "" {
				status = "ACTIVE"
			}
			promo := entities.Promotion{
				Name: req.Name, Description: req.Description, Type: req.Type,
				PlayHours: req.PlayHours, FreeHours: req.FreeHours,
				DiscountPct: req.DiscountPct, DiscountAmt: req.DiscountAmt,
				TableTypes: req.TableTypes, StartDate: startDate, EndDate: endDate,
				Status: status, CreatedBy: ctx.GetString("UserId"),
			}
			result, err := repository.Promotion.CreatePromotion(promo)
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusCreated, result)
		})

	r.PUT("/:promotionId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("promotionId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_001, "invalid promotionId")
				return
			}
			var req request.Promotion
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_001, err.Error())
				return
			}
			startDate, _ := time.Parse("2006-01-02", req.StartDate)
			endDate, _ := time.Parse("2006-01-02", req.EndDate)
			promo := entities.Promotion{
				Name: req.Name, Description: req.Description, Type: req.Type,
				PlayHours: req.PlayHours, FreeHours: req.FreeHours,
				DiscountPct: req.DiscountPct, DiscountAmt: req.DiscountAmt,
				TableTypes: req.TableTypes, StartDate: startDate, EndDate: endDate,
				Status: req.Status, UpdatedBy: ctx.GetString("UserId"),
			}
			if err := repository.Promotion.UpdatePromotionById(id, promo); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	r.DELETE("/:promotionId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("promotionId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_001, "invalid promotionId")
				return
			}
			if err := repository.Promotion.DeletePromotionById(id); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.PM_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})
}
