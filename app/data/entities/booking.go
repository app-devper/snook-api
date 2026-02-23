package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	TableId       primitive.ObjectID `bson:"tableId" json:"tableId"`
	TableName     string             `bson:"tableName" json:"tableName"`
	CustomerName  string             `bson:"customerName" json:"customerName"`
	CustomerPhone string             `bson:"customerPhone" json:"customerPhone"`
	BookingDate   time.Time          `bson:"bookingDate" json:"bookingDate"`
	StartTime     string             `bson:"startTime" json:"startTime"`
	EndTime       string             `bson:"endTime" json:"endTime"`
	Status        string             `bson:"status" json:"status"`
	Note          string             `bson:"note" json:"note"`
	CreatedBy     string             `bson:"createdBy" json:"-"`
	CreatedDate   time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy     string             `bson:"updatedBy" json:"-"`
	UpdatedDate   time.Time          `bson:"updatedDate" json:"-"`
}
