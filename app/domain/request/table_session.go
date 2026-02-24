package request

type OpenTable struct {
	TableId string `json:"tableId" binding:"required"`
}

type CloseTable struct {
	TableCharge *float64 `json:"tableCharge"`
	Discount    float64  `json:"discount"`
	Note        string   `json:"note"`
	PaymentType string   `json:"paymentType"`
	PaymentNote string   `json:"paymentNote"`
}

type TransferTable struct {
	NewTableId string `json:"newTableId" binding:"required"`
}

type ApplyPromotion struct {
	PromotionId string `json:"promotionId" binding:"required"`
}
