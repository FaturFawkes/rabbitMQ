package main

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err.Error())
	}
	log.Println("Success Call RabbitMQ!")
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to conncect RabbitMQ")
	defer conn.Close()	

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", 	// name
		false, 		// durable
		false, 		// delete when unused
		false, 		// exclusive
		false, 		// no-wait
		nil, 		// arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, 	// queue
		"", 		// consumer
		true, 		// auto-ack
		false, 		// exclusive
		false, 		// no-local
		false, 		// no-wait
		nil, 		// arguments
	)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			time.Sleep(1 * time.Second)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}