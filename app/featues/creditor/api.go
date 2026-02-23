package creditor

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/domain"
	"snook/app/domain/request"
	"snook/middlewares"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyCreditorAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("creditors")

	r.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		status := ctx.Query("status")
		creditors, err := repository.Creditor.GetCreditors(status)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.CR_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, creditors)
	})

	r.GET("/:creditorId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("creditorId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_001, "invalid creditorId")
			return
		}
		creditorRecord, err := repository.Creditor.GetCreditorById(id)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, creditorRecord)
	})

	r.GET("/:creditorId/payments", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("creditorId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_001, "invalid creditorId")
			return
		}
		payments, err := repository.Creditor.GetCreditorPayments(id)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.CR_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, payments)
	})

	r.POST("/:creditorId/pay", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("creditorId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_001, "invalid creditorId")
			return
		}
		var req request.CreditorPayment
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_001, err.Error())
			return
		}
		creditorRecord, err := repository.Creditor.GetCreditorById(id)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_002, "creditor not found")
			return
		}
		if req.Amount <= 0 {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_001, "amount must be greater than 0")
			return
		}
		if req.Amount > creditorRecord.Remaining {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_001, "amount exceeds remaining balance")
			return
		}
		payment := entities.CreditorPayment{
			CreditorId: id, Amount: req.Amount, Type: req.Type,
			Note: req.Note, CreatedBy: ctx.GetString("UserId"),
		}
		_, err = repository.Creditor.CreateCreditorPayment(payment)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.CR_BAD_REQUEST_002, err.Error())
			return
		}
		creditorRecord.PaidAmount += req.Amount
		creditorRecord.Remaining = creditorRecord.Amount - creditorRecord.PaidAmount
		if creditorRecord.Remaining <= 0 {
			creditorRecord.Remaining = 0
			creditorRecord.Status = "PAID"
		}
		creditorRecord.UpdatedBy = ctx.GetString("UserId")
		_ = repository.Creditor.UpdateCreditor(id, creditorRecord)
		ctx.JSON(http.StatusOK, creditorRecord)
	})
}
