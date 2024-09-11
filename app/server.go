package main

import (
	"go-gin/app/handlers"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Middleware to log requests
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, duration)
	})

	// Route to handle POST orders
	r.POST("/orders", handlers.OrderHandler)

	// Route to handle GET all orders
	r.GET("/orders", handlers.GetAllOrdersHandler)

	// // Route to handle GET a specific order by ID
	r.GET("/orders/:id", handlers.GetOrderHandler)

	// // Route to handle PUT to update a specific order by ID
	r.PUT("/orders/:id", handlers.UpdateOrderHandler)

	// // Route to handle DELETE a specific order by ID
	r.DELETE("/orders/:id", handlers.DeleteOrderHandler)

	// Start the Gin server
	r.Run(":8080")
}
