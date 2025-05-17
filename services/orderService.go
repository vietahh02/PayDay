package services

import (
	"errors"
	"os"
	"payday/initializers"
	"payday/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CreateOrder(order *models.OrderModel) (models.OrderModel, error) {
	result := initializers.DB.Create(&order)
	if result.Error != nil {
		return models.OrderModel{}, result.Error
	}
	return *order, nil
}

func UpdateOrder(transactionId, orderID, status, bankCode, message string) (models.OrderModel, error) {

	var order models.OrderModel
	err := initializers.DB.Where("transaction_id = ? AND order_id = ?", transactionId, orderID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.OrderModel{}, errors.New("order not found")
		} else {
			return models.OrderModel{}, errors.New("system error")
		}
	}
	order.Status = status
	order.BankCode = bankCode
	order.Message = message

	result := initializers.DB.Save(&order)
	if result.Error != nil {
		return models.OrderModel{}, result.Error
	}
	return order, nil
}

func GenerateJWT() (string, error) {
	t := time.Now().Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"iss":     os.Getenv("PARTNER_CODE"),
		"api_key": os.Getenv("API_KEY"),
		"jti":     os.Getenv("API_KEY") + strconv.Itoa(int(t)),
		"exp":     t,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetOrderByTransactionId(transactionId string) (models.OrderModel, error) {
	var order models.OrderModel
	err := initializers.DB.Where("transaction_id = ?", transactionId).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.OrderModel{}, errors.New("order not found")
		} else {
			return models.OrderModel{}, errors.New("system error")
		}
	}
	return order, nil
}
