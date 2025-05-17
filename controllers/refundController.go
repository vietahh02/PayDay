package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"payday/models"
	"payday/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateRefund(c *gin.Context) {

	type Body struct {
		TransactionID string `json:"transactionId"`
		Reason        string `json:"reason"`
	}

	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error-1": err.Error(),
		})
		return
	}

	type RefundRequest struct {
		PartnerRefID  string `json:"partnerRefId"`
		TransactionID string `json:"transactionId"`
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
		Reason        string `json:"reason"`
	}

	order, err := services.GetOrderByTransactionId(body.TransactionID)
	if err != nil {
		c.JSON(400, gin.H{
			"error-2": err.Error(),
		})
		return
	}

	refund := RefundRequest{
		PartnerRefID:  uuid.NewString(),
		TransactionID: body.TransactionID,
		Amount:        order.Amount,
		Currency:      order.Currency,
		Reason:        body.Reason,
	}

	jwt, err := services.GenerateJWT()
	if err != nil {
		c.JSON(400, gin.H{
			"error0": err.Error(),
		})
		return
	}
	b, _ := json.Marshal(refund)
	httpReq, err := http.NewRequest("POST", os.Getenv("URL_REFUND"), bytes.NewReader(b))
	httpReq.Header.Set("X-APPOTAPAY-AUTH", jwt)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Request-ID", uuid.NewString())
	httpReq.Header.Set("X-Language", "vi")
	if err != nil {
		c.JSON(400, gin.H{
			"error-2": err.Error(),
		})
		return
	}

	client := http.Client{}
	resp, err := client.Do(httpReq)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.JSON(400, gin.H{
				"error4": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorResponse,
		})
	} else {
		var refundResponse models.RefundModel
		if err := json.NewDecoder(resp.Body).Decode(&refundResponse); err != nil {
			c.JSON(400, gin.H{
				"error3": err.Error(),
			})
			return
		}
		_, err = services.CreateRefund(&refundResponse)
		if err != nil {
			c.JSON(400, gin.H{
				"error5": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"response": refundResponse,
		})
	}

}

type ErrorResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Errors    []struct {
		Field  string `json:"field"`
		Reason string `json:"reason"`
	} `json:"errors"`
}
