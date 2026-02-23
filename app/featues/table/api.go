package table

import (
	"snook/app/core/constant"
	"snook/app/domain"
	"snook/app/featues/table/usecase"
	"snook/middlewares"

	"github.com/gin-gonic/gin"
)

func ApplyTableAPI(
	route *gin.RouterGroup,
	repository *domain.Repository,
) {
	tableRoute := route.Group("tables")

	tableRoute.GET("",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.GetTables(repository.Table),
	)

	tableRoute.POST("",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN),
		usecase.CreateTable(repository.Table),
	)

	tableRoute.PUT("/:tableId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN),
		usecase.UpdateTableById(repository.Table),
	)

	tableRoute.PATCH("/:tableId/status",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN),
		usecase.UpdateTableStatus(repository.Table),
	)

	tableRoute.DELETE("/:tableId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN),
		usecase.DeleteTableById(repository.Table),
	)
}
