package rabbitmq

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func StartConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"orders", // queue name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			processOrder(string(msg.Body))

			if err := processOrder(string(msg.Body)); err != nil {
				log.Printf("Failed to process order: %v", err)
				// Implement retry logic or send to dead letter exchange
			}
		}
	}()

	// Log that the consumer has started and is waiting for messages
	log.Println("Consumer started. Waiting for messages...")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down consumer...")
}

func processOrder(orderID string) error {
	// Log the order being processed
	log.Printf("Processing order: %s", orderID)
	// Simulate order processing (e.g., saving to DB, external API call)

	// Return an error if processing fails
	return nil
}
