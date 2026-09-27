package table_session

import (
	"snook/app/domain"
	"snook/app/featues/table_session/usecase"
	"snook/middlewares"

	"github.com/gin-gonic/gin"
)

func ApplyTableSessionAPI(
	route *gin.RouterGroup,
	repository *domain.Repository,
) {
	sessionRoute := route.Group("sessions")

	sessionRoute.GET("",
		middlewares.RequireSession(repository.Auth),
		usecase.GetTableSessions(repository.TableSession),
	)

	sessionRoute.GET("/:sessionId",
		middlewares.RequireSession(repository.Auth),
		usecase.GetTableSessionById(repository.TableSession, repository.TableOrder, repository.Payment),
	)

	sessionRoute.GET("/table/:tableId/active",
		middlewares.RequireSession(repository.Auth),
		usecase.GetActiveSessionByTableId(repository.TableSession),
	)

	sessionRoute.POST("/open",
		middlewares.RequireSession(repository.Auth),
		usecase.OpenTable(repository.TableSession, repository.Table),
	)

	sessionRoute.POST("/:sessionId/close",
		middlewares.RequireSession(repository.Auth),
		usecase.CloseTable(repository.TableSession, repository.Table, repository.TableOrder, repository.Payment, repository.Promotion),
	)

	sessionRoute.POST("/:sessionId/pause",
		middlewares.RequireSession(repository.Auth),
		usecase.PauseTable(repository.TableSession),
	)

	sessionRoute.POST("/:sessionId/resume",
		middlewares.RequireSession(repository.Auth),
		usecase.ResumeTable(repository.TableSession),
	)

	sessionRoute.POST("/:sessionId/transfer",
		middlewares.RequireSession(repository.Auth),
		usecase.TransferTable(repository.TableSession, repository.Table),
	)

	sessionRoute.POST("/:sessionId/apply-promotion",
		middlewares.RequireSession(repository.Auth),
		usecase.ApplyPromotionToSession(repository.TableSession, repository.Promotion),
	)
}
