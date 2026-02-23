package request

type Expense struct {
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required"`
}
