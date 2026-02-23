package menu

import (
	"net/http"
	"snook/app/core/constant"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/domain"
	"snook/app/domain/request"
	"snook/middlewares"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyMenuAPI(route *gin.RouterGroup, repository *domain.Repository) {
	// ─── Categories ─────────────────────────────────
	catRoute := route.Group("menu-categories")

	catRoute.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		cats, err := repository.MenuCategory.GetMenuCategories()
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.MC_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, cats)
	})

	catRoute.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			var req request.MenuCategory
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_001, err.Error())
				return
			}
			cat := entities.MenuCategory{Name: req.Name, SortOrder: req.SortOrder, CreatedBy: ctx.GetString("UserId")}
			result, err := repository.MenuCategory.CreateMenuCategory(cat)
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusCreated, result)
		})

	catRoute.PUT("/:categoryId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("categoryId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_001, "invalid categoryId")
				return
			}
			var req request.MenuCategory
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_001, err.Error())
				return
			}
			cat := entities.MenuCategory{Name: req.Name, SortOrder: req.SortOrder, UpdatedBy: ctx.GetString("UserId")}
			if err := repository.MenuCategory.UpdateMenuCategoryById(id, cat); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	catRoute.DELETE("/:categoryId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("categoryId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_001, "invalid categoryId")
				return
			}
			if err := repository.MenuCategory.DeleteMenuCategoryById(id); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MC_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	// ─── Menu Items ─────────────────────────────────
	itemRoute := route.Group("menu-items")

	itemRoute.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		category := ctx.Query("category")
		items, err := repository.MenuItem.GetMenuItems(category)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.MI_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, items)
	})

	itemRoute.GET("/:itemId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("itemId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, "invalid itemId")
			return
		}
		item, err := repository.MenuItem.GetMenuItemById(id)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, item)
	})

	itemRoute.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			var req request.MenuItem
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, err.Error())
				return
			}
			status := req.Status
			if status == "" {
				status = "ACTIVE"
			}
			item := entities.MenuItem{
				Name: req.Name, Category: req.Category, Price: req.Price,
				CostPrice: req.CostPrice, Quantity: req.Quantity, Unit: req.Unit,
				Status: status, ImageUrl: req.ImageUrl, CreatedBy: ctx.GetString("UserId"),
			}
			result, err := repository.MenuItem.CreateMenuItem(item)
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusCreated, result)
		})

	itemRoute.PUT("/:itemId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("itemId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, "invalid itemId")
				return
			}
			var req request.MenuItem
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, err.Error())
				return
			}
			item := entities.MenuItem{
				Name: req.Name, Category: req.Category, Price: req.Price,
				CostPrice: req.CostPrice, Quantity: req.Quantity, Unit: req.Unit,
				Status: req.Status, ImageUrl: req.ImageUrl, UpdatedBy: ctx.GetString("UserId"),
			}
			if err := repository.MenuItem.UpdateMenuItemById(id, item); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	itemRoute.DELETE("/:itemId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("itemId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, "invalid itemId")
				return
			}
			if err := repository.MenuItem.DeleteMenuItemById(id); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	itemRoute.PATCH("/:itemId/quantity", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session),
		middlewares.RequireAuthorization(constant.SUPER, constant.ADMIN), func(ctx *gin.Context) {
			id, err := primitive.ObjectIDFromHex(ctx.Param("itemId"))
			if err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, "invalid itemId")
				return
			}
			var req request.MenuItemQuantity
			if err := ctx.ShouldBindJSON(&req); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_001, err.Error())
				return
			}
			if err := repository.MenuItem.UpdateMenuItemQuantity(id, req.Quantity); err != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.MI_BAD_REQUEST_002, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	itemRoute.GET("/low-stock", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		threshold := 10
		if t := ctx.Query("threshold"); t != "" {
			if v, err := strconv.Atoi(t); err == nil {
				threshold = v
			}
		}
		items, err := repository.MenuItem.GetLowStockMenuItems(threshold)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.MI_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, items)
	})
}
