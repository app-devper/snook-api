package report

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/domain"
	"snook/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyReportAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("reports")

	r.GET("/revenue", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		startDate := ctx.Query("startDate")
		endDate := ctx.Query("endDate")
		if startDate == "" || endDate == "" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.RP_BAD_REQUEST_001, "startDate and endDate required")
			return
		}
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24*time.Hour - time.Nanosecond)

		sessions, err := repository.TableSession.GetTableSessions(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.RP_INTERNAL_001, err.Error())
			return
		}
		expenses, err := repository.Expense.GetExpenses(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.RP_INTERNAL_001, err.Error())
			return
		}

		totalIncome := 0.0
		totalTableCharge := 0.0
		totalFoodIncome := 0.0
		for _, s := range sessions {
			if s.Status == "CLOSED" {
				totalIncome += s.GrandTotal
				totalTableCharge += s.TableCharge
				totalFoodIncome += s.FoodTotal
			}
		}
		totalExpense := 0.0
		for _, e := range expenses {
			totalExpense += e.Amount
		}

		ctx.JSON(http.StatusOK, gin.H{
			"startDate":        startDate,
			"endDate":          endDate,
			"totalSessions":    len(sessions),
			"totalIncome":      totalIncome,
			"totalTableCharge": totalTableCharge,
			"totalFoodIncome":  totalFoodIncome,
			"totalExpense":     totalExpense,
			"netProfit":        totalIncome - totalExpense,
			"sessions":         sessions,
			"expenses":         expenses,
		})
	})

	r.GET("/revenue/by-table/:tableId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		tableId, err := primitive.ObjectIDFromHex(ctx.Param("tableId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.RP_BAD_REQUEST_001, "invalid tableId")
			return
		}
		startDate := ctx.Query("startDate")
		endDate := ctx.Query("endDate")
		if startDate == "" || endDate == "" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.RP_BAD_REQUEST_001, "startDate and endDate required")
			return
		}
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.RP_BAD_REQUEST_001, "invalid startDate format")
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.RP_BAD_REQUEST_001, "invalid endDate format")
			return
		}
		end = end.Add(24*time.Hour - time.Nanosecond)
		sessions, err := repository.TableSession.GetSessionsByTableId(tableId, start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.RP_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, sessions)
	})
}
