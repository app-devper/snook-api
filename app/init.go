package app

import (
	"os"
	"snook/app/domain"
	"snook/app/featues/booking"
	"snook/app/featues/creditor"
	"snook/app/featues/dashboard"
	"snook/app/featues/expense"
	"snook/app/featues/menu"
	"snook/app/featues/payment"
	"snook/app/featues/promotion"
	"snook/app/featues/report"
	"snook/app/featues/setting"
	"snook/app/featues/table"
	"snook/app/featues/table_order"
	"snook/app/featues/table_session"
	"snook/db"
	"snook/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Routes struct{}

func (app Routes) StartGin() {
	r := gin.New()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		logrus.Error(err)
	}

	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())
	r.Use(middlewares.NewCors([]string{"*"}))

	resource, err := db.InitResource()
	if err != nil {
		logrus.Fatal("failed to init database: ", err)
	}
	defer resource.Close()

	publicRoute := r.Group("/api/snook/v1")

	repository := domain.InitRepository(resource)

	table.ApplyTableAPI(publicRoute, repository)
	table_session.ApplyTableSessionAPI(publicRoute, repository)
	booking.ApplyBookingAPI(publicRoute, repository)
	menu.ApplyMenuAPI(publicRoute, repository)
	table_order.ApplyTableOrderAPI(publicRoute, repository)
	payment.ApplyPaymentAPI(publicRoute, repository)
	creditor.ApplyCreditorAPI(publicRoute, repository)
	promotion.ApplyPromotionAPI(publicRoute, repository)
	expense.ApplyExpenseAPI(publicRoute, repository)
	setting.ApplySettingAPI(publicRoute, repository)
	dashboard.ApplyDashboardAPI(publicRoute, repository)
	report.ApplyReportAPI(publicRoute, repository)

	r.NoRoute(middlewares.NoRoute())

	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		logrus.Error(err)
	}
}
