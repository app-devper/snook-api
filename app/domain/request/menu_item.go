package request

type MenuItem struct {
	Name      string  `json:"name" binding:"required"`
	Category  string  `json:"category" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	CostPrice float64 `json:"costPrice"`
	Quantity  int     `json:"quantity"`
	Unit      string  `json:"unit"`
	Status    string  `json:"status"`
	ImageUrl  string  `json:"imageUrl"`
}

type MenuItemQuantity struct {
	Quantity int `json:"quantity" binding:"required"`
}
