package models

type Order struct {
	Transaction struct {
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
		BankCode      string `json:"bankCode"`
		PaymentMethod string `json:"paymentMethod"`
		Action        string `json:"action"`
	} `json:"transaction"`
	PartnerReference struct {
		Order struct {
			ID        string `json:"id"`
			Info      string `json:"info"`
			ExtraData string `json:"extraData"`
		} `json:"order"`
		NotificationConfig struct {
			NotifyURL   string `json:"notifyUrl"`
			RedirectURL string `json:"redirectUrl"`
		} `json:"notificationConfig"`
	} `json:"partnerReference"`
}
