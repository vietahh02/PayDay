package interfaces

import "payday/models"

type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
}
