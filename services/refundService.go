package services

import (
	"payday/initializers"
	"payday/models"
)

func CreateRefund(refund *models.RefundModel) (models.RefundModel, error) {
	result := initializers.DB.Create(&refund)
	if result.Error != nil {
		return models.RefundModel{}, result.Error
	}
	return *refund, nil
}
