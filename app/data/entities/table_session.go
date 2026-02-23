package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableSession struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	TableId         primitive.ObjectID `bson:"tableId" json:"tableId"`
	TableName       string             `bson:"tableName" json:"tableName"`
	TableType       string             `bson:"tableType" json:"tableType"`
	RatePerHour     float64            `bson:"ratePerHour" json:"ratePerHour"`
	Status          string             `bson:"status" json:"status"`
	StartTime       time.Time          `bson:"startTime" json:"startTime"`
	EndTime         *time.Time         `bson:"endTime,omitempty" json:"endTime,omitempty"`
	PausedAt        *time.Time         `bson:"pausedAt,omitempty" json:"pausedAt,omitempty"`
	TotalPausedMins float64            `bson:"totalPausedMins" json:"totalPausedMins"`
	DurationMins    float64            `bson:"durationMins" json:"durationMins"`
	TableCharge     float64            `bson:"tableCharge" json:"tableCharge"`
	FoodTotal       float64            `bson:"foodTotal" json:"foodTotal"`
	Discount        float64            `bson:"discount" json:"discount"`
	PromotionId     *primitive.ObjectID `bson:"promotionId,omitempty" json:"promotionId,omitempty"`
	PromotionName   string             `bson:"promotionName" json:"promotionName"`
	PromotionDiscount float64          `bson:"promotionDiscount" json:"promotionDiscount"`
	GrandTotal      float64            `bson:"grandTotal" json:"grandTotal"`
	Note            string             `bson:"note" json:"note"`
	CreatedBy       string             `bson:"createdBy" json:"-"`
	CreatedDate     time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy       string             `bson:"updatedBy" json:"-"`
	UpdatedDate     time.Time          `bson:"updatedDate" json:"-"`
}

type TableSessionDetail struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	TableId         primitive.ObjectID `bson:"tableId" json:"tableId"`
	TableName       string             `bson:"tableName" json:"tableName"`
	TableType       string             `bson:"tableType" json:"tableType"`
	RatePerHour     float64            `bson:"ratePerHour" json:"ratePerHour"`
	Status          string             `bson:"status" json:"status"`
	StartTime       time.Time          `bson:"startTime" json:"startTime"`
	EndTime         *time.Time         `bson:"endTime,omitempty" json:"endTime,omitempty"`
	PausedAt        *time.Time         `bson:"pausedAt,omitempty" json:"pausedAt,omitempty"`
	TotalPausedMins float64            `bson:"totalPausedMins" json:"totalPausedMins"`
	DurationMins    float64            `bson:"durationMins" json:"durationMins"`
	TableCharge     float64            `bson:"tableCharge" json:"tableCharge"`
	FoodTotal       float64            `bson:"foodTotal" json:"foodTotal"`
	Discount        float64            `bson:"discount" json:"discount"`
	PromotionId     *primitive.ObjectID `bson:"promotionId,omitempty" json:"promotionId,omitempty"`
	PromotionName   string             `bson:"promotionName" json:"promotionName"`
	PromotionDiscount float64          `bson:"promotionDiscount" json:"promotionDiscount"`
	GrandTotal      float64            `bson:"grandTotal" json:"grandTotal"`
	Note            string             `bson:"note" json:"note"`
	CreatedBy       string             `bson:"createdBy" json:"-"`
	CreatedDate     time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy       string             `bson:"updatedBy" json:"-"`
	UpdatedDate     time.Time          `bson:"updatedDate" json:"-"`
	Orders          []TableOrder       `json:"orders"`
	Payments        []Payment          `json:"payments"`
}

type SessionSummary struct {
	TotalSessions int     `bson:"totalSessions" json:"totalSessions"`
	TotalRevenue  float64 `bson:"totalRevenue" json:"totalRevenue"`
	TotalTable    float64 `bson:"totalTable" json:"totalTable"`
	TotalFood     float64 `bson:"totalFood" json:"totalFood"`
}

type SessionDailyChart struct {
	Date          string  `bson:"_id" json:"date"`
	TotalSessions int     `bson:"totalSessions" json:"totalSessions"`
	TotalRevenue  float64 `bson:"totalRevenue" json:"totalRevenue"`
}
