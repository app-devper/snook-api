package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuItem struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Price       float64            `bson:"price" json:"price"`
	CostPrice   float64            `bson:"costPrice" json:"costPrice"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Unit        string             `bson:"unit" json:"unit"`
	Status      string             `bson:"status" json:"status"`
	ImageUrl    string             `bson:"imageUrl" json:"imageUrl"`
	CreatedBy   string             `bson:"createdBy" json:"-"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	UpdatedBy   string             `bson:"updatedBy" json:"-"`
	UpdatedDate time.Time          `bson:"updatedDate" json:"-"`
}

type LowStockMenuItem struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Category string             `bson:"category" json:"category"`
	Quantity int                `bson:"quantity" json:"quantity"`
	Unit     string             `bson:"unit" json:"unit"`
}
