package models

type Order2 struct {
	Transaction struct {
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
		BankCode      string `json:"bankCode"`
		PaymentMethod string `json:"paymentMethod"`
		Action        string `json:"action"`
	} `json:"transaction"`

	SourceOfFunds struct {
		Type string `json:"type"`
		Card struct {
			CardNumber     string `json:"cardNumber"`
			CardHolderName string `json:"cardHolderName"`
			CardMonth      string `json:"cardMonth"`
			CardYear       string `json:"cardYear"`
		} `json:"card"`
	} `json:"sourceOfFunds"`

	PartnerReference struct {
		Order struct {
			ID        string `json:"id"`
			Info      string `json:"info"`
			ExtraData string `json:"extraData"`
		} `json:"order"`
		NotificationConfig struct {
			NotifyUrl   string `json:"notifyUrl"`
			RedirectUrl string `json:"redirectUrl"`
		} `json:"notificationConfig"`
	} `json:"partnerReference"`
}
