package request

type Table struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	RatePerHour float64 `json:"ratePerHour" binding:"required"`
	Description string  `json:"description"`
}

type TableStatus struct {
	Status string `json:"status" binding:"required"`
}
