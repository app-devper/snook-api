package request

type OpenTable struct {
	TableId string `json:"tableId" binding:"required"`
}

type CloseTable struct {
	Discount float64 `json:"discount"`
	Note     string  `json:"note"`
}

type TransferTable struct {
	NewTableId string `json:"newTableId" binding:"required"`
}

type ApplyPromotion struct {
	PromotionId string `json:"promotionId" binding:"required"`
}
