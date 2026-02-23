package usecase

import (
	"math"
	"net/http"
	"snook/app/core/errcode"
	"snook/app/data/entities"
	"snook/app/data/repositories"
	"snook/app/domain/request"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OpenTable(sessionEntity repositories.ITableSession, tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req request.OpenTable
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, err.Error())
			return
		}
		tableId, err := primitive.ObjectIDFromHex(req.TableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid tableId")
			return
		}
		table, err := tableEntity.GetTableById(tableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "table not found")
			return
		}
		if table.Status != "AVAILABLE" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "table is not available")
			return
		}
		userId := ctx.GetString("UserId")
		session := entities.TableSession{
			TableId:     tableId,
			TableName:   table.Name,
			TableType:   table.Type,
			RatePerHour: table.RatePerHour,
			Status:      "ACTIVE",
			StartTime:   time.Now(),
			CreatedBy:   userId,
		}
		result, err := sessionEntity.CreateTableSession(session)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, err.Error())
			return
		}
		_ = tableEntity.UpdateTableStatus(tableId, "IN_USE")
		ctx.JSON(http.StatusCreated, result)
	}
}

func CloseTable(sessionEntity repositories.ITableSession, tableEntity repositories.ITable, orderEntity repositories.ITableOrder, paymentEntity repositories.IPayment, promotionEntity repositories.IPromotion) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		session, err := sessionEntity.GetTableSessionById(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session not found")
			return
		}
		if session.Status == "CLOSED" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session already closed")
			return
		}
		var req request.CloseTable
		if err := ctx.ShouldBindJSON(&req); err == nil {
			session.Discount = req.Discount
			session.Note = req.Note
		}
		now := time.Now()
		session.EndTime = &now
		totalMins := now.Sub(session.StartTime).Minutes() - session.TotalPausedMins
		if totalMins < 0 {
			totalMins = 0
		}
		session.DurationMins = math.Round(totalMins*100) / 100
		billableMins := totalMins
		if billableMins < 60 {
			billableMins = 60
		}
		session.TableCharge = math.Round((billableMins/60)*session.RatePerHour*100) / 100

		// Recalculate promotion discount at close time
		if session.PromotionId != nil {
			promo, promoErr := promotionEntity.GetPromotionById(*session.PromotionId)
			if promoErr == nil {
				promoDiscount := 0.0
				switch promo.Type {
				case "FREE_HOURS":
					if promo.PlayHours > 0 && (totalMins/60) >= promo.PlayHours {
						promoDiscount = promo.FreeHours * session.RatePerHour
					}
				case "DISCOUNT_PCT":
					promoDiscount = math.Round(session.TableCharge*promo.DiscountPct) / 100
				case "DISCOUNT_AMT":
					promoDiscount = promo.DiscountAmt
				}
				session.PromotionDiscount = math.Round(promoDiscount*100) / 100
			}
		}

		orders, _ := orderEntity.GetOrdersBySessionId(sessionId)
		foodTotal := 0.0
		for _, o := range orders {
			foodTotal += o.Total
		}
		session.FoodTotal = foodTotal
		session.GrandTotal = session.TableCharge + session.FoodTotal - session.Discount - session.PromotionDiscount
		if session.GrandTotal < 0 {
			session.GrandTotal = 0
		}

		// Sum existing payments
		payments, _ := paymentEntity.GetPaymentsBySessionId(sessionId)
		paidTotal := 0.0
		for _, p := range payments {
			paidTotal += p.Amount
		}

		// Auto-create final payment for remaining balance
		remaining := math.Round((session.GrandTotal-paidTotal)*100) / 100
		if remaining > 0 {
			payType := req.PaymentType
			if payType == "" {
				payType = "CASH"
			}
			userId := ctx.GetString("UserId")
			_, payErr := paymentEntity.CreatePayment(entities.Payment{
				SessionId:   sessionId,
				Type:        payType,
				Amount:      remaining,
				Note:        req.PaymentNote,
				CreatedBy:   userId,
				CreatedDate: now,
			})
			if payErr != nil {
				errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "failed to create payment: "+payErr.Error())
				return
			}
		}

		session.Status = "CLOSED"
		userId := ctx.GetString("UserId")
		session.UpdatedBy = userId
		if err := sessionEntity.UpdateTableSession(sessionId, session); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, err.Error())
			return
		}
		_ = tableEntity.UpdateTableStatus(session.TableId, "AVAILABLE")
		ctx.JSON(http.StatusOK, session)
	}
}

func PauseTable(sessionEntity repositories.ITableSession) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		session, err := sessionEntity.GetTableSessionById(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session not found")
			return
		}
		if session.Status != "ACTIVE" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session is not active")
			return
		}
		now := time.Now()
		session.PausedAt = &now
		session.Status = "PAUSED"
		if err := sessionEntity.UpdateTableSession(sessionId, session); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, session)
	}
}

func ResumeTable(sessionEntity repositories.ITableSession) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		session, err := sessionEntity.GetTableSessionById(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session not found")
			return
		}
		if session.Status != "PAUSED" || session.PausedAt == nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session is not paused")
			return
		}
		pausedDuration := time.Since(*session.PausedAt).Minutes()
		session.TotalPausedMins += pausedDuration
		session.PausedAt = nil
		session.Status = "ACTIVE"
		if err := sessionEntity.UpdateTableSession(sessionId, session); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, session)
	}
}

func TransferTable(sessionEntity repositories.ITableSession, tableEntity repositories.ITable) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		var req request.TransferTable
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, err.Error())
			return
		}
		newTableId, err := primitive.ObjectIDFromHex(req.NewTableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid newTableId")
			return
		}
		session, err := sessionEntity.GetTableSessionById(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session not found")
			return
		}
		newTable, err := tableEntity.GetTableById(newTableId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "new table not found")
			return
		}
		if newTable.Status != "AVAILABLE" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "new table is not available")
			return
		}
		oldTableId := session.TableId
		session.TableId = newTableId
		session.TableName = newTable.Name
		session.TableType = newTable.Type
		session.RatePerHour = newTable.RatePerHour
		userId := ctx.GetString("UserId")
		session.UpdatedBy = userId
		if err := sessionEntity.UpdateTableSession(sessionId, session); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, err.Error())
			return
		}
		_ = tableEntity.UpdateTableStatus(oldTableId, "AVAILABLE")
		_ = tableEntity.UpdateTableStatus(newTableId, "IN_USE")
		ctx.JSON(http.StatusOK, session)
	}
}

func ApplyPromotionToSession(sessionEntity repositories.ITableSession, promotionEntity repositories.IPromotion) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		var req request.ApplyPromotion
		if err := ctx.ShouldBindJSON(&req); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, err.Error())
			return
		}
		promotionId, err := primitive.ObjectIDFromHex(req.PromotionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid promotionId")
			return
		}
		session, err := sessionEntity.GetTableSessionById(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session not found")
			return
		}
		if session.Status == "CLOSED" {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session already closed")
			return
		}
		promo, err := promotionEntity.GetPromotionById(promotionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "promotion not found")
			return
		}
		// Calculate promotion discount based on type
		promoDiscount := 0.0
		elapsed := time.Since(session.StartTime).Minutes() - session.TotalPausedMins
		if elapsed < 0 {
			elapsed = 0
		}
		switch promo.Type {
		case "FREE_HOURS":
			if promo.PlayHours > 0 && (elapsed/60) >= promo.PlayHours {
				promoDiscount = promo.FreeHours * session.RatePerHour
			}
		case "DISCOUNT_PCT":
			currentCharge := (elapsed / 60) * session.RatePerHour
			promoDiscount = math.Round(currentCharge*promo.DiscountPct) / 100
		case "DISCOUNT_AMT":
			promoDiscount = promo.DiscountAmt
		}
		session.PromotionId = &promotionId
		session.PromotionName = promo.Name
		session.PromotionDiscount = math.Round(promoDiscount*100) / 100
		session.UpdatedBy = ctx.GetString("UserId")
		if err := sessionEntity.UpdateTableSession(sessionId, session); err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, session)
	}
}

func GetTableSessions(sessionEntity repositories.ITableSession) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startDate := ctx.Query("startDate")
		endDate := ctx.Query("endDate")
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid startDate")
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid endDate")
			return
		}
		end = end.Add(24*time.Hour - time.Nanosecond)
		sessions, err := sessionEntity.GetTableSessions(start, end)
		if err != nil {
			errcode.Abort(ctx, http.StatusInternalServerError, errcode.TS_INTERNAL_001, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, sessions)
	}
}

func GetTableSessionById(sessionEntity repositories.ITableSession, orderEntity repositories.ITableOrder, paymentEntity repositories.IPayment) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := primitive.ObjectIDFromHex(ctx.Param("sessionId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid sessionId")
			return
		}
		session, err := sessionEntity.GetTableSessionById(sessionId)
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_002, "session not found")
			return
		}
		orders, _ := orderEntity.GetOrdersBySessionId(sessionId)
		payments, _ := paymentEntity.GetPaymentsBySessionId(sessionId)
		detail := entities.TableSessionDetail{
			Id: session.Id, TableId: session.TableId, TableName: session.TableName,
			TableType: session.TableType, RatePerHour: session.RatePerHour,
			Status: session.Status, StartTime: session.StartTime, EndTime: session.EndTime,
			PausedAt: session.PausedAt, TotalPausedMins: session.TotalPausedMins,
			DurationMins: session.DurationMins, TableCharge: session.TableCharge,
			FoodTotal: session.FoodTotal, Discount: session.Discount,
			PromotionId: session.PromotionId, PromotionName: session.PromotionName,
			PromotionDiscount: session.PromotionDiscount, GrandTotal: session.GrandTotal,
			Note: session.Note, CreatedDate: session.CreatedDate,
			Orders: orders, Payments: payments,
		}
		ctx.JSON(http.StatusOK, detail)
	}
}

func GetActiveSessionByTableId(sessionEntity repositories.ITableSession) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tableId, err := primitive.ObjectIDFromHex(ctx.Param("tableId"))
		if err != nil {
			errcode.Abort(ctx, http.StatusBadRequest, errcode.TS_BAD_REQUEST_001, "invalid tableId")
			return
		}
		session, err := sessionEntity.GetActiveSessionByTableId(tableId)
		if err != nil {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusOK, session)
	}
}
