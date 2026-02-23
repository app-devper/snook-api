package dashboard

import (
	"net/http"
	"snook/app/core/errcode"
	"snook/app/domain"
	"snook/middlewares"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ApplyDashboardAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("dashboard")

	r.GET("/summary", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		startDate := ctx.DefaultQuery("startDate", time.Now().Format("2006-01-02"))
		endDate := ctx.DefaultQuery("endDate", time.Now().Format("2006-01-02"))
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24*time.Hour - time.Nanosecond)
		summary, err := repository.TableSession.GetSessionSummary(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.DA_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, summary)
	})

	r.GET("/daily-chart", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		startDate := ctx.DefaultQuery("startDate", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
		endDate := ctx.DefaultQuery("endDate", time.Now().Format("2006-01-02"))
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24*time.Hour - time.Nanosecond)
		chart, err := repository.TableSession.GetSessionDailyChart(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.DA_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, chart)
	})

	r.GET("/low-stock", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		threshold := 10
		if t := ctx.Query("threshold"); t != "" {
			if v, err := strconv.Atoi(t); err == nil {
				threshold = v
			}
		}
		items, err := repository.MenuItem.GetLowStockMenuItems(threshold)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.DA_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, items)
	})
}
