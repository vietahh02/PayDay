package models

import (
	"gorm.io/gorm"
)

type OrderModel struct {
	gorm.Model
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	BankCode      string `json:"bankCode"`
	Action        string `json:"action"`
	OrderID       string `json:"orderId"`
	Info          string `json:"info"`
	ExtraData     string `json:"extraData"`
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	URL           string `json:"url"`
}
