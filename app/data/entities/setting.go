package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	CompanyName    string             `bson:"companyName" json:"companyName"`
	CompanyAddress string             `bson:"companyAddress" json:"companyAddress"`
	CompanyPhone   string             `bson:"companyPhone" json:"companyPhone"`
	CompanyTaxId   string             `bson:"companyTaxId" json:"companyTaxId"`
	ReceiptFooter  string             `bson:"receiptFooter" json:"receiptFooter"`
	PromptPayId    string             `bson:"promptPayId" json:"promptPayId"`
	UpdatedBy      string             `bson:"updatedBy" json:"-"`
	UpdatedDate    time.Time          `bson:"updatedDate" json:"-"`
}
