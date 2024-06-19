package utils

type PostPaymentInput struct {
	CardNumber  string `json:"cardnumber" binding:"required" validate:"card_number" `
	ExpiryDate  string `json:"expirydate" binding:"required" validate:"expiry_date" `
	Cvv         string `json:"cvv" validate:"cvv" binding:"required"`
	Amount      string `json:"amount" binding:"required" validate:"amount"`
	Description string `json:"description" `
}