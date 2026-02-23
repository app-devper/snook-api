package request

type CreditorPayment struct {
	Amount float64 `json:"amount" binding:"required"`
	Type   string  `json:"type" binding:"required"`
	Note   string  `json:"note"`
}
