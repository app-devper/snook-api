package table

import (
	"snook/app/domain"
	"snook/app/featues/table/usecase"
	"snook/middlewares"

	"github.com/app-devper/um-api/sessionclient"
	"github.com/gin-gonic/gin"
)

func ApplyTableAPI(
	route *gin.RouterGroup,
	repository *domain.Repository,
) {
	tableRoute := route.Group("tables")

	tableRoute.GET("",
		middlewares.RequireSession(repository.Auth),
		usecase.GetTables(repository.Table),
	)

	tableRoute.POST("",
		middlewares.RequireSession(repository.Auth),
		repository.Auth.AtLeast(sessionclient.RoleAdmin),
		usecase.CreateTable(repository.Table),
	)

	tableRoute.PUT("/:tableId",
		middlewares.RequireSession(repository.Auth),
		repository.Auth.AtLeast(sessionclient.RoleAdmin),
		usecase.UpdateTableById(repository.Table),
	)

	tableRoute.PATCH("/:tableId/status",
		middlewares.RequireSession(repository.Auth),
		repository.Auth.AtLeast(sessionclient.RoleAdmin),
		usecase.UpdateTableStatus(repository.Table),
	)

	tableRoute.DELETE("/:tableId",
		middlewares.RequireSession(repository.Auth),
		repository.Auth.AtLeast(sessionclient.RoleAdmin),
		usecase.DeleteTableById(repository.Table),
	)
}
