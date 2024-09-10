package main

import (
	"go-gin/app/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Route to handle POST orders
	r.POST("/orders", handlers.OrderHandler)

	// Route to handle GET all orders
	r.GET("/orders", handlers.GetAllOrdersHandler)

	// // Route to handle GET a specific order by ID
	// r.GET("/orders/:id", handlers.GetOrderHandler)

	// // Route to handle PUT to update a specific order by ID
	// r.PUT("/orders/:id", handlers.UpdateOrderHandler)

	// // Route to handle DELETE a specific order by ID
	// r.DELETE("/orders/:id", handlers.DeleteOrderHandler)

	// Start the Gin server
	r.Run(":8080")
}
