package request

type Booking struct {
	TableId       string `json:"tableId" binding:"required"`
	CustomerName  string `json:"customerName" binding:"required"`
	CustomerPhone string `json:"customerPhone"`
	BookingDate   string `json:"bookingDate" binding:"required"`
	StartTime     string `json:"startTime" binding:"required"`
	EndTime       string `json:"endTime" binding:"required"`
	Note          string `json:"note"`
}

type BookingStatus struct {
	Status string `json:"status" binding:"required"`
}
