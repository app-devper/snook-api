package request

type Setting struct {
	CompanyName    string `json:"companyName"`
	CompanyAddress string `json:"companyAddress"`
	CompanyPhone   string `json:"companyPhone"`
	CompanyTaxId   string `json:"companyTaxId"`
	ReceiptFooter  string `json:"receiptFooter"`
	PromptPayId    string `json:"promptPayId"`
}
