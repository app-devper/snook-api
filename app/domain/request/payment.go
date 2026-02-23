package request

type Payment struct {
	SessionId string  `json:"sessionId" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Note      string  `json:"note"`
}

type Checkout struct {
	Payments []PaymentItem `json:"payments" binding:"required"`
	Discount float64       `json:"discount"`
	Note     string        `json:"note"`
}

type PaymentItem struct {
	Type   string  `json:"type" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Note   string  `json:"note"`
}
