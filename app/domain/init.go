package domain

import (
	"snook/app/data/repositories"
	"snook/db"
)

type Repository struct {
	Session      repositories.ISession
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
		Session:      repositories.NewSessionEntity(resource),
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
