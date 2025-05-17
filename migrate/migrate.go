package main

import (
	"payday/initializers"
	"payday/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.OrderModel{})
	initializers.DB.AutoMigrate(&models.RefundModel{})
}
