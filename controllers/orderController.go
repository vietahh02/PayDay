package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"payday/initializers"
	"payday/models"
	"payday/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePayment(c *gin.Context) {
	type Req struct {
		Amount        int    `json:"amount" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
		Info          string `json:"info" `
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error-1": err.Error(),
		})
		return
	}

	var order models.Order
	order.Transaction.Amount = req.Amount
	order.Transaction.Currency = "VND"
	order.Transaction.BankCode = ""
	order.Transaction.PaymentMethod = req.PaymentMethod
	order.Transaction.Action = "PAY"
	order.PartnerReference.Order.ID = uuid.New().String()
	order.PartnerReference.Order.Info = req.Info
	order.PartnerReference.Order.ExtraData = "Extra data"
	order.PartnerReference.NotificationConfig.NotifyURL = os.Getenv("URL_FREE_CLOUD_FLARE") + "/receive_notify"
	order.PartnerReference.NotificationConfig.RedirectURL = os.Getenv("URL_FREE_CLOUD_FLARE") + "/redirect"

	// Call AppotaPay
	client := http.Client{Timeout: 10 * time.Second}
	jwt, err := services.GenerateJWT()
	if err != nil {
		c.JSON(400, gin.H{
			"error0": err.Error(),
		})
		return
	}
	body, _ := json.Marshal(order)
	httpReq, err := http.NewRequest("POST", os.Getenv("URL_PAYMENT_NORMAL"), bytes.NewReader(body))
	httpReq.Header.Set("X-APPOTAPAY-AUTH", jwt)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Request-ID", uuid.NewString())
	httpReq.Header.Set("X-Language", "vi")

	if err != nil {
		c.JSON(400, gin.H{
			"error1": err.Error(),
		})
		return
	}
	resp, err := client.Do(httpReq)

	if err != nil || resp.StatusCode != 200 {
		c.JSON(400, gin.H{
			"error2": "Data input wrong",
		})
		return
	}
	defer resp.Body.Close()
	var apResp struct {
		Transaction struct {
			TransactionID string `json:"transactionId"`
		} `json:"transaction"`
		Payment struct {
			URL string `json:"url"`
		} `json:"payment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apResp); err != nil {
		c.JSON(400, gin.H{
			"error3": err.Error(),
		})
		return
	}

	orderModel := &models.OrderModel{
		Amount:        order.Transaction.Amount,
		Currency:      order.Transaction.Currency,
		BankCode:      order.Transaction.BankCode,
		Action:        order.Transaction.Action,
		OrderID:       order.PartnerReference.Order.ID,
		Info:          order.PartnerReference.Order.Info,
		ExtraData:     order.PartnerReference.Order.ExtraData,
		TransactionID: apResp.Transaction.TransactionID,
		Status:        "pending",
		URL:           apResp.Payment.URL,
	}

	_, err = services.CreateOrder(orderModel)
	if err != nil {
		c.JSON(400, gin.H{
			"error4": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_url":    apResp.Payment.URL,
		"transaction_id": apResp.Transaction.TransactionID,
	})
}

func CreatePaymentBank(c *gin.Context) {
	type Req struct {
		Amount int    `json:"amount" binding:"required"`
		Info   string `json:"info" `
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error-1": err.Error(),
		})
		return
	}

	var order models.Order2
	order.Transaction.Amount = req.Amount
	order.Transaction.Currency = "VND"
	order.Transaction.BankCode = "VPBANK"
	order.Transaction.PaymentMethod = "ATM"
	order.Transaction.Action = "PAY"
	order.SourceOfFunds.Type = "card"
	order.SourceOfFunds.Card.CardNumber = "9704000000000018"
	order.SourceOfFunds.Card.CardHolderName = "Nguyen Van A"
	order.SourceOfFunds.Card.CardMonth = "03"
	order.SourceOfFunds.Card.CardYear = "07"
	order.PartnerReference.Order.ID = uuid.New().String()
	order.PartnerReference.Order.Info = req.Info
	order.PartnerReference.Order.ExtraData = "Extra data"
	order.PartnerReference.NotificationConfig.NotifyUrl = os.Getenv("URL_FREE_CLOUD_FLARE") + "/receive_notify"
	order.PartnerReference.NotificationConfig.RedirectUrl = os.Getenv("URL_FREE_CLOUD_FLARE") + "/redirect"

	// Call AppotaPay
	client := http.Client{Timeout: 10 * time.Second}
	jwt, err := services.GenerateJWT()
	if err != nil {
		c.JSON(400, gin.H{
			"error0": err.Error(),
		})
		return
	}
	body, _ := json.Marshal(order)
	httpReq, err := http.NewRequest("POST", os.Getenv("URL_PAYMENT_BANK"), bytes.NewReader(body))
	httpReq.Header.Set("X-APPOTAPAY-AUTH", jwt)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Request-ID", uuid.NewString())
	httpReq.Header.Set("X-Language", "vi")

	if err != nil {
		c.JSON(400, gin.H{
			"error1": err.Error(),
		})
		return
	}
	resp, err := client.Do(httpReq)

	if err != nil || resp.StatusCode != 200 {
		c.JSON(400, gin.H{
			"error2": "Data input wrong",
		})
		return
	}
	defer resp.Body.Close()
	var apResp struct {
		Transaction struct {
			TransactionID string `json:"transactionId"`
			Status        string `json:"status"`
			ErrorCode     int    `json:"errorCode"`
			ErrorMessage  string `json:"errorMessage"`
			PartnerCode   string `json:"partnerCode"`
			Amount        int    `json:"amount"`
			OrderAmount   int    `json:"orderAmount"`
			Fee           struct {
				CustomerFee int `json:"customer_fee"`
			} `json:"fee"`
			Currency      string `json:"currency"`
			BankCode      string `json:"bankCode"`
			PaymentMethod string `json:"paymentMethod"`
			Action        string `json:"action"`
			CreatedAt     string `json:"createdAt"`
			UpdatedAt     string `json:"updatedAt"`
		} `json:"transaction"`

		Authentication struct {
			VerificationUrl    string `json:"verificationUrl"`
			VerificationMethod string `json:"verificationMethod"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"authentication"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apResp); err != nil {
		c.JSON(400, gin.H{
			"error3": err.Error(),
		})
		return
	}

	orderModel := &models.OrderModel{
		Amount:        order.Transaction.Amount,
		Currency:      order.Transaction.Currency,
		BankCode:      order.Transaction.BankCode,
		Action:        order.Transaction.Action,
		OrderID:       order.PartnerReference.Order.ID,
		Info:          order.PartnerReference.Order.Info,
		ExtraData:     order.PartnerReference.Order.ExtraData,
		TransactionID: apResp.Transaction.TransactionID,
		Status:        "pending",
		URL:           apResp.Authentication.VerificationUrl,
	}

	_, err = services.CreateOrder(orderModel)
	if err != nil {
		c.JSON(400, gin.H{
			"error4": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment_url":    apResp.Authentication.VerificationUrl,
		"transaction_id": apResp.Transaction.TransactionID,
	})
}

func ReceiveNotifyPayment(c *gin.Context) {

	type Notify struct {
		Data      string `json:"data"`
		Signature string `json:"signature"`
	}

	var notify Notify
	if err := c.ShouldBindJSON(&notify); err != nil {
		c.JSON(400, gin.H{
			"error-1": err.Error(),
		})
		return
	}
	jsonBytes, err := base64.StdEncoding.DecodeString(notify.Data)
	if err != nil {
		c.JSON(400, gin.H{
			"error1": err.Error(),
		})
	}

	var result ResponseData
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		c.JSON(400, gin.H{
			"error2": err.Error(),
		})
	}

	_, err = services.UpdateOrder(result.Transaction.TransactionID, result.PartnerReference.Order.ID, result.Transaction.Status, result.Transaction.BankCode, result.Transaction.ErrorMessage)
	if err != nil {
		c.JSON(400, gin.H{
			"error-2": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func RedirectPayment(c *gin.Context) {
	data := c.Query("data")

	jsonBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		c.JSON(400, gin.H{
			"error1": err.Error(),
		})
	}
	var result ResponseData
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		c.JSON(400, gin.H{
			"error2": err.Error(),
		})
	}
	if result.Transaction.Status != "success" {
		_, err = services.UpdateOrder(result.Transaction.TransactionID, result.PartnerReference.Order.ID, result.Transaction.Status, result.Transaction.BankCode, result.Transaction.ErrorMessage)
		if err != nil {
			c.JSON(400, gin.H{
				"error-2": err.Error(),
			})
		}
	}

	c.JSON(200, gin.H{
		"message": result,
	})
}

func CheckOrder(c *gin.Context) {

	transactionId := c.Param("transactionId")

	var order models.OrderModel
	err := initializers.DB.Where("transaction_id = ?", transactionId).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"order": order,
	})

}

type Transaction struct {
	TransactionID    string `json:"transactionId"`
	ReconciliationID string `json:"reconciliationId"`
	PartnerCode      string `json:"partnerCode"`
	Status           string `json:"status"`
	ErrorCode        int    `json:"errorCode"`
	ErrorMessage     string `json:"errorMessage"`
	OrderAmount      int    `json:"orderAmount"`
	Amount           int    `json:"amount"`
	DiscountAmount   int    `json:"discountAmount"`
	Currency         string `json:"currency"`
	BankCode         string `json:"bankCode"`
	PaymentMethod    string `json:"paymentMethod"`
	Action           string `json:"action"`
	ClientIP         string `json:"clientIp"`
	Version          string `json:"version"`
	Fee              Fee    `json:"fee"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type Fee struct {
	CustomerFee int `json:"customer_fee"`
}

type OrderInfo struct {
	ID        string `json:"id"`
	Info      string `json:"info"`
	ExtraData string `json:"extraData"`
}

type PartnerReference struct {
	Order OrderInfo `json:"order"`
}

type ResponseData struct {
	Transaction      Transaction      `json:"transaction"`
	PartnerReference PartnerReference `json:"partnerReference"`
}
