package models

import "gorm.io/gorm"

type RefundModel struct {
	gorm.Model
	RefundID      string `json:"refundId"`
	TransactionID string `json:"transactionId"`
	PartnerRefID  string `json:"partnerRefId"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Reason        string `json:"reason"`
	Status        string `json:"status"`
	MessageError  string `json:"messageError"`
	RefundedAt    string `json:"refundedAt"`
}
