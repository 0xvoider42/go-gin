# Go microservice with RabbitMQ and Gin
This is a simple Go microservice that uses the Gin web framework for HTTP handling and RabbitMQ for messaging. It includes features such as dead-letter queues (DLQ) and supports integration with RabbitMQ for message processing.

## Features
- RESTful API: Exposes an HTTP API to receive and queue orders.
- RabbitMQ Integration: Publishes and consumes messages to/from RabbitMQ.
- Dead-Letter Queue (DLQ): Implements a DLQ to handle rejected messages.
- Mockable Components: RabbitMQ interactions are mockable for testing purposes.

## Running this thing

1. Install go and you are ready to go
2. `go mod tidy` to install dependencies
3. `docker compose up` to build and run docker image for rabbit
4. `go run app/sever.go` to run the app

## Usage
- `{
    "ID": "5", "item": "just a cammel", "price": 1200, "message_type": "success"
}` example of a request
- `(POST) http://localhost:8080/orders` to post orders
- `(GET) http://localhost:8080/orders` to get all orders
- `(PUT) http://localhost:8080/orders/:id` to update order
- `(DELETE) http://localhost:8080/orders/:id` to delete order

### Extra
at the moment there is inline documentation
