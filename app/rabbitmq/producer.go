package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func ConnectRabbitMQ() (*amqp.Channel, *amqp.Connection, error) {
	var err error

	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(
		"orders_dlq", // DLQ name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("failed to declare the orders DLQ: %v", err)
		return nil, nil, fmt.Errorf("failed to declare the orders queue: %v", err)
	}

	_, err = ch.QueueInspect("orders")
	if err != nil {
		_, err = ch.QueueDeclare(
			"orders", // queue name
			true,     // durable
			false,    // delete when unused
			false,    // exclusive
			false,    // no-wait
			amqp.Table{
				"x-dead-letter-exchange":    "",
				"x-dead-letter-routing-key": "orders_dlq",
			},
		)
		if err != nil {
			log.Printf("failed to declare the orders DLQ: %v", err)
			return nil, nil, fmt.Errorf("failed to declare the orders queue: %v", err)
		}
	} else {
		_, err = ch.QueueInspect("orders_dlq")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to inspect the orders DLQ: %v", err)
		}
	}

	// Declare a topic exchange
	err = ch.ExchangeDeclare(
		"order_topic", // exchange name
		"direct",      // exchange type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Bind queue to the exchange with the "order" routing key
	err = ch.QueueBind(
		"orders",      // queue name
		"orders",      // routing key
		"order_topic", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to bind a queue: %v", err)
	}

	return ch, conn, nil
}

func PublishMessage(ch *amqp.Channel, message string, messageType string) error {
	// Construct the routing key
	routingKey := "order." + messageType

	// Publish the message to the exchange
	err := ch.Publish(
		"order_topic", // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}
	log.Printf("Message sent: %s with routing key: %s", message, routingKey)
	return nil
}
