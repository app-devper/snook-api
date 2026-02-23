package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Creditor struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	SessionId     primitive.ObjectID `bson:"sessionId" json:"sessionId"`
	CustomerName  string             `bson:"customerName" json:"customerName"`
	CustomerPhone string             `bson:"customerPhone" json:"customerPhone"`
	Amount        float64            `bson:"amount" json:"amount"`
	PaidAmount    float64            `bson:"paidAmount" json:"paidAmount"`
	Remaining     float64            `bson:"remaining" json:"remaining"`
	Status        string             `bson:"status" json:"status"`
	Note          string             `bson:"note" json:"note"`
	DueDate       *time.Time         `bson:"dueDate,omitempty" json:"dueDate,omitempty"`
	CreatedBy     string             `bson:"createdBy" json:"-"`
	CreatedDate   time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy     string             `bson:"updatedBy" json:"-"`
	UpdatedDate   time.Time          `bson:"updatedDate" json:"-"`
}

type CreditorPayment struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	CreditorId primitive.ObjectID `bson:"creditorId" json:"creditorId"`
	Amount     float64            `bson:"amount" json:"amount"`
	Type       string             `bson:"type" json:"type"`
	Note       string             `bson:"note" json:"note"`
	CreatedBy  string             `bson:"createdBy" json:"-"`
	CreatedDate time.Time         `bson:"createdDate" json:"createdDate"`
}
