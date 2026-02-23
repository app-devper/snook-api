package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	SessionId   primitive.ObjectID `bson:"sessionId" json:"sessionId"`
	Type        string             `bson:"type" json:"type"`
	Amount      float64            `bson:"amount" json:"amount"`
	Note        string             `bson:"note" json:"note"`
	CreatedBy   string             `bson:"createdBy" json:"-"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
}
