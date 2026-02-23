package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableOrder struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	SessionId  primitive.ObjectID `bson:"sessionId" json:"sessionId"`
	MenuItemId primitive.ObjectID `bson:"menuItemId" json:"menuItemId"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	CostPrice  float64            `bson:"costPrice" json:"costPrice"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Discount   float64            `bson:"discount" json:"discount"`
	Total      float64            `bson:"total" json:"total"`
	CreatedBy  string             `bson:"createdBy" json:"-"`
	CreatedDate time.Time         `bson:"createdDate" json:"createdDate"`
}
