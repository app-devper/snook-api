package payment

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/domain"
	"snook/app/domain/request"
	"snook/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyPaymentAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("payments")

	r.GET("/session/:sessionId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PY_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		payments, err := repository.Payment.GetPaymentsBySessionId(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.PY_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, payments)
	})

	r.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		startDate := ctx.DefaultQuery("startDate", time.Now().Format("2006-01-02"))
		endDate := ctx.DefaultQuery("endDate", time.Now().Format("2006-01-02"))
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24*time.Hour - time.Nanosecond)
		payments, err := repository.Payment.GetPaymentsByDateRange(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.PY_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, payments)
	})

	r.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		var req request.Payment
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PY_BAD_REQUEST_001, err.Error())
			return
		}
		sessionId, _ := primitive.ObjectIDFromHex(req.SessionId)
		payment := entities.Payment{
			SessionId: sessionId, Type: req.Type, Amount: req.Amount,
			Note: req.Note, CreatedBy: ctx.GetString("UserId"),
		}
		// If type is OUTSTANDING, create a creditor record
		if req.Type == "OUTSTANDING" {
			session, err := repository.TableSession.GetTableSessionById(sessionId)
			if err == nil {
				creditor := entities.Creditor{
					SessionId: sessionId, CustomerName: req.Note,
					Amount: req.Amount, Remaining: req.Amount,
					Status: "PENDING", CreatedBy: ctx.GetString("UserId"),
				}
				_, _ = repository.Creditor.CreateCreditor(creditor)
				_ = session // used for context
			}
		}
		result, err := repository.Payment.CreatePayment(payment)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PY_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, result)
	})

	r.DELETE("/:paymentId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("paymentId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PY_BAD_REQUEST_001, "invalid paymentId")
			return
		}
		if err := repository.Payment.DeletePayment(id); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.PY_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}
