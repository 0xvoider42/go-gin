package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go-gin/app/rabbitmq"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID          string `json:"id"`
	Item        string `json:"item"`
	Price       int    `json:"price"`
	MessageType string `json:"message_type"`
}

func OrderHandler(c *gin.Context) {
	var order Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

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
	// Connect to RabbitMQ
	ch, conn, err := rabbitmq.ConnectRabbitMQ()

	log.Println("1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to RabbitMQ"})
		return
	}
	defer conn.Close()
	defer ch.Close()

	log.Println("2")

	// Declare the queue (assuming the queue name is "orders")
	queue, err := ch.QueueDeclare(
		"orders", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to declare queue"})
		return
	}

	log.Println("3")

	// Consume messages from the queue
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to consume messages"})
		return
	}

	// Create a slice to hold the orders
	var orders []Order

	// Loop through the messages and append to orders slice
	for msg := range msgs {
		var order Order
		if err := json.Unmarshal(msg.Body, &order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse order message"})
			return
		}
		orders = append(orders, order)
	}

	log.Println(orders)

	// Return the list of orders
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func GetOrderHandler(c *gin.Context) {
	// Implement logic to get a specific order by ID
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Get order", "orderID": orderID})
}

func UpdateOrderHandler(c *gin.Context) {
	// Implement logic to update a specific order by ID
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Update order", "orderID": orderID})
}

func DeleteOrderHandler(c *gin.Context) {
	// Implement logic to delete a specific order by ID
	orderID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Delete order", "orderID": orderID})
}
