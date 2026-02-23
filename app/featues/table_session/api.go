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
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.GetTableSessions(repository.TableSession),
	)

	sessionRoute.GET("/:sessionId",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.GetTableSessionById(repository.TableSession, repository.TableOrder, repository.Payment),
	)

	sessionRoute.GET("/table/:tableId/active",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.GetActiveSessionByTableId(repository.TableSession),
	)

	sessionRoute.POST("/open",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.OpenTable(repository.TableSession, repository.Table),
	)

	sessionRoute.POST("/:sessionId/close",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.CloseTable(repository.TableSession, repository.Table, repository.TableOrder, repository.Payment, repository.Promotion),
	)

	sessionRoute.POST("/:sessionId/pause",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.PauseTable(repository.TableSession),
	)

	sessionRoute.POST("/:sessionId/resume",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.ResumeTable(repository.TableSession),
	)

	sessionRoute.POST("/:sessionId/transfer",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.TransferTable(repository.TableSession, repository.Table),
	)

	sessionRoute.POST("/:sessionId/apply-promotion",
		middlewares.RequireAuthenticated(),
		middlewares.RequireSession(repository.Session),
		usecase.ApplyPromotionToSession(repository.TableSession, repository.Promotion),
	)
}
