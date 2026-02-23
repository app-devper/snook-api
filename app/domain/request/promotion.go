package request

type Promotion struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Type        string   `json:"type" binding:"required"`
	PlayHours   float64  `json:"playHours"`
	FreeHours   float64  `json:"freeHours"`
	DiscountPct float64  `json:"discountPct"`
	DiscountAmt float64  `json:"discountAmt"`
	TableTypes  []string `json:"tableTypes"`
	StartDate   string   `json:"startDate" binding:"required"`
	EndDate     string   `json:"endDate" binding:"required"`
	Status      string   `json:"status"`
}
