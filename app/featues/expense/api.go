package expense

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

func ApplyExpenseAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("expenses")

	r.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		startDate := ctx.DefaultQuery("startDate", time.Now().Format("2006-01-02"))
		endDate := ctx.DefaultQuery("endDate", time.Now().Format("2006-01-02"))
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24*time.Hour - time.Nanosecond)
		expenses, err := repository.Expense.GetExpenses(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.EX_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, expenses)
	})

	r.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		var req request.Expense
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_001, err.Error())
			return
		}
		date, _ := time.Parse("2006-01-02", req.Date)
		expense := entities.Expense{
			Category: req.Category, Description: req.Description,
			Amount: req.Amount, Date: date, CreatedBy: ctx.GetString("UserId"),
		}
		result, err := repository.Expense.CreateExpense(expense)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, result)
	})

	r.PUT("/:expenseId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("expenseId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_001, "invalid expenseId")
			return
		}
		var req request.Expense
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_001, err.Error())
			return
		}
		date, _ := time.Parse("2006-01-02", req.Date)
		expense := entities.Expense{
			Category: req.Category, Description: req.Description,
			Amount: req.Amount, Date: date, UpdatedBy: ctx.GetString("UserId"),
		}
		if err := repository.Expense.UpdateExpenseById(id, expense); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.DELETE("/:expenseId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("expenseId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_001, "invalid expenseId")
				return
			}
			if err := repository.Expense.DeleteExpenseById(id); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.EX_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})
}
