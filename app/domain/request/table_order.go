package request

type TableOrder struct {
	SessionId  string  `json:"sessionId" binding:"required"`
	MenuItemId string  `json:"menuItemId" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required"`
	Discount   float64 `json:"discount"`
}

type UpdateTableOrder struct {
	Quantity int     `json:"quantity" binding:"required"`
	Discount float64 `json:"discount"`
}
