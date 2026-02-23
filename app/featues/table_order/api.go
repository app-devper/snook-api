package table_order

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

func ApplyTableOrderAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("table-orders")

	r.GET("/session/:sessionId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		orders, err := repository.TableOrder.GetOrdersBySessionId(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.TO_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, orders)
	})

	r.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		var req request.TableOrder
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_001, err.Error())
			return
		}
		sessionId, _ := primitive.ObjectIDFromHex(req.SessionId)
		menuItemId, _ := primitive.ObjectIDFromHex(req.MenuItemId)
		menuItem, err := repository.MenuItem.GetMenuItemById(menuItemId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_002, "menu item not found")
			return
		}
		if menuItem.Quantity < req.Quantity {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_002, "insufficient stock")
			return
		}
		total := (menuItem.Price * float64(req.Quantity)) - req.Discount
		if total < 0 {
			total = 0
		}
		order := entities.TableOrder{
			SessionId: sessionId, MenuItemId: menuItemId,
			Name: menuItem.Name, Price: menuItem.Price, CostPrice: menuItem.CostPrice,
			Quantity: req.Quantity, Discount: req.Discount, Total: total,
			CreatedBy: ctx.GetString("UserId"),
		}
		result, err := repository.TableOrder.CreateTableOrder(order)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_002, err.Error())
			return
		}
		_ = repository.MenuItem.UpdateMenuItemQuantity(menuItemId, -req.Quantity)
		ctx.JSON(http.StatusCreated, result)
	})

	r.DELETE("/:orderId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		orderId, err := primitive.ObjectIDFromHex(ctx.Param("orderId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_001, "invalid orderId")
			return
		}
		order, err := repository.TableOrder.GetTableOrderById(orderId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_002, "order not found")
			return
		}
		if err := repository.TableOrder.DeleteTableOrder(orderId); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TO_BAD_REQUEST_002, err.Error())
			return
		}
		_ = repository.MenuItem.UpdateMenuItemQuantity(order.MenuItemId, order.Quantity)
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}
