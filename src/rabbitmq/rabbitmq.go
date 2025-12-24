package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

func InitRabbitMQ() error {
	var err error

	// Подключение к RabbitMQ с ретраями
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/5): %v", i+1, err)
		time.Sleep(time.Second * 3)
	}

	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		"email_notifications",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	return err
}

func Close() {
	if ch != nil {
		_ = ch.Close()
	}
	if conn != nil {
		_ = conn.Close()
	}
}

func SendSuccessfulMessage(message string) error {
	if ch == nil {
		if err := InitRabbitMQ(); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx,
		"",                    // exchange
		"email_notifications", // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Сохранять сообщения при перезапуске
			ContentType:  "text/plain",
			Body:         []byte(message),
			Timestamp:    time.Now(),
		})
}
