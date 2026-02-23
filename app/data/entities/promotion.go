package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Promotion struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Type        string             `bson:"type" json:"type"`
	PlayHours   float64            `bson:"playHours" json:"playHours"`
	FreeHours   float64            `bson:"freeHours" json:"freeHours"`
	DiscountPct float64            `bson:"discountPct" json:"discountPct"`
	DiscountAmt float64            `bson:"discountAmt" json:"discountAmt"`
	TableTypes  []string           `bson:"tableTypes" json:"tableTypes"`
	StartDate   time.Time          `bson:"startDate" json:"startDate"`
	EndDate     time.Time          `bson:"endDate" json:"endDate"`
	Status      string             `bson:"status" json:"status"`
	CreatedBy   string             `bson:"createdBy" json:"-"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"-"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"-"`
}
