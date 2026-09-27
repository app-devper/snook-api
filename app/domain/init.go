package domain

import (
	"github.com/app-devper/um-api/sessionclient/ginauth"
	"snook/app/data/repositories"
	"snook/db"
)

type Repository struct {
	// Auth verifies UM tokens and sessions; set by the app at startup.
	Auth         *ginauth.Auth
	Table        repositories.ITable
	TableSession repositories.ITableSession
	Booking      repositories.IBooking
	MenuCategory repositories.IMenuCategory
	MenuItem     repositories.IMenuItem
	TableOrder   repositories.ITableOrder
	Payment      repositories.IPayment
	Creditor     repositories.ICreditor
	Promotion    repositories.IPromotion
	Expense      repositories.IExpense
	Setting      repositories.ISetting
}

func InitRepository(resource *db.Resource) *Repository {
	return &Repository{
		Table:        repositories.NewTableEntity(resource),
		TableSession: repositories.NewTableSessionEntity(resource),
		Booking:      repositories.NewBookingEntity(resource),
		MenuCategory: repositories.NewMenuCategoryEntity(resource),
		MenuItem:     repositories.NewMenuItemEntity(resource),
		TableOrder:   repositories.NewTableOrderEntity(resource),
		Payment:      repositories.NewPaymentEntity(resource),
		Creditor:     repositories.NewCreditorEntity(resource),
		Promotion:    repositories.NewPromotionEntity(resource),
		Expense:      repositories.NewExpenseEntity(resource),
		Setting:      repositories.NewSettingEntity(resource),
	}
}
