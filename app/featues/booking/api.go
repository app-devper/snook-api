package booking

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

func ApplyBookingAPI(route *gin.RouterGroup, repository *domain.Repository) {
	r := route.Group("bookings")

	r.GET("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		startDate := ctx.DefaultQuery("startDate", time.Now().Format("2006-01-02"))
		endDate := ctx.DefaultQuery("endDate", time.Now().Format("2006-01-02"))
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		end = end.Add(24*time.Hour - time.Nanosecond)
		bookings, err := repository.Booking.GetBookings(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.BK_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, bookings)
	})

	r.GET("/:bookingId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("bookingId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, "invalid bookingId")
			return
		}
		booking, err := repository.Booking.GetBookingById(id)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, booking)
	})

	r.POST("", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		var req request.Booking
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, err.Error())
			return
		}
		tableId, _ := primitive.ObjectIDFromHex(req.TableId)
		table, err := repository.Table.GetTableById(tableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, "table not found")
			return
		}
		bookingDate, _ := time.Parse("2006-01-02", req.BookingDate)
		userId := ctx.GetString("UserId")
		booking := entities.Booking{
			TableId: tableId, TableName: table.Name,
			CustomerName: req.CustomerName, CustomerPhone: req.CustomerPhone,
			BookingDate: bookingDate, StartTime: req.StartTime, EndTime: req.EndTime,
			Status: "PENDING", Note: req.Note, CreatedBy: userId,
		}
		result, err := repository.Booking.CreateBooking(booking)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, result)
	})

	r.PUT("/:bookingId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("bookingId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, "invalid bookingId")
			return
		}
		var req request.Booking
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, err.Error())
			return
		}
		tableId, err := primitive.ObjectIDFromHex(req.TableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, "invalid tableId")
			return
		}
		table, err := repository.Table.GetTableById(tableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, "table not found")
			return
		}
		bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, "invalid bookingDate format")
			return
		}
		userId := ctx.GetString("UserId")
		booking := entities.Booking{
			TableId: tableId, TableName: table.Name,
			CustomerName: req.CustomerName, CustomerPhone: req.CustomerPhone,
			BookingDate: bookingDate, StartTime: req.StartTime, EndTime: req.EndTime,
			Note: req.Note, UpdatedBy: userId,
		}
		if err := repository.Booking.UpdateBookingById(id, booking); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.PATCH("/:bookingId/status", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("bookingId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, "invalid bookingId")
			return
		}
		var req request.BookingStatus
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, err.Error())
			return
		}
		if err := repository.Booking.UpdateBookingStatus(id, req.Status); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.DELETE("/:bookingId", middlewares.RequireAuthenticated(), middlewares.RequireSession(repository.Session), func(ctx *gin.Context) {
		id, err := primitive.ObjectIDFromHex(ctx.Param("bookingId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_001, "invalid bookingId")
			return
		}
		if err := repository.Booking.DeleteBookingById(id); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.BK_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}
