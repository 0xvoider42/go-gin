package handlers

import (
	"log"
	"net/http"
	"sync"

	"go-gin/app/rabbitmq"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID          string `json:"id"`
	Item        string `json:"item"`
	Price       int    `json:"price"`
	MessageType string `json:"message_type"`
}

// Global map to store orders
var orders = make(map[string]Order)
var mu sync.Mutex

func OrderHandler(c *gin.Context) {
	var order Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Save order to the global map
	mu.Lock()
	orders[order.ID] = order
	mu.Unlock()

	ch, conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to RabbitMQ"})
		return
	}
	// Ensure the connection and channel are closed when the function exits
	defer conn.Close()
	defer ch.Close()

	// Send order to RabbitMQ
	err = rabbitmq.PublishMessage(ch, order.ID, order.MessageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue the order"})
		return
	}

	// Log the successful receipt and queuing of the order
	log.Printf("Order received and queued: ID=%s, MessageType=%s", order.ID, order.MessageType)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Order received and queued",
		"orderID":     order.ID,
		"messageType": order.MessageType,
	})
}

func GetAllOrdersHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// GetOrderHandler handles fetching a specific order by ID
func GetOrderHandler(c *gin.Context) {
	orderID := c.Param("id")

	mu.Lock()
	order, exists := orders[orderID]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// UpdateOrderHandler handles updating a specific order by ID
func UpdateOrderHandler(c *gin.Context) {
	orderID := c.Param("id")

	var updatedOrder Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	mu.Lock()
	order, exists := orders[orderID]
	if exists {
		// Update the order details
		order.Item = updatedOrder.Item
		order.Price = updatedOrder.Price
		order.MessageType = updatedOrder.MessageType
		orders[orderID] = order
	}
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ch, conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to RabbitMQ"})
		return
	}
	defer conn.Close()
	defer ch.Close()

	err = rabbitmq.PublishMessage(ch, order.ID, "update")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue the order update"})
		return
	}

	log.Printf("Order updated and queued: ID=%s, MessageType=%s", order.ID, "update")

	c.JSON(http.StatusOK, gin.H{"message": "Order updated", "order": order})
}

// DeleteOrderHandler handles deleting a specific order by ID
func DeleteOrderHandler(c *gin.Context) {
	orderID := c.Param("id")

	mu.Lock()
	_, exists := orders[orderID]
	if exists {
		delete(orders, orderID)
	}
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ch, conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to RabbitMQ"})
		return
	}
	defer conn.Close()
	defer ch.Close()

	err = rabbitmq.PublishMessage(ch, orderID, "delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue the order deletion"})
		return
	}

	log.Printf("Order deleted and queued: ID=%s, MessageType=%s", orderID, "delete")

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted", "orderID": orderID})
}
