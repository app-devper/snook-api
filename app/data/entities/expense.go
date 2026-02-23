package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Amount      float64            `bson:"amount" json:"amount"`
	Date        time.Time          `bson:"date" json:"date"`
	CreatedBy   string             `bson:"createdBy" json:"-"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"-"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"-"`
}
