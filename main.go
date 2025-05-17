package main

import (
	"payday/controllers"
	"payday/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "run oke nhó",
		})
	})

	//order
	router.POST("receive_notify", controllers.ReceiveNotifyPayment)
	router.GET("redirect", controllers.RedirectPayment)
	router.POST("/payment", controllers.CreatePayment)
	router.POST("/payment_bank", controllers.CreatePaymentBank)
	router.GET("/check_order/:transactionId", controllers.CheckOrder)

	//refund
	router.POST("/refund/", controllers.CreateRefund)

	router.Run()
}
