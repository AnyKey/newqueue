package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs := getQueue(ch)

	forever := make(chan bool)

	file, err := os.OpenFile("/home/redhat/docker-elk/logs/main.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	go func() {
		for d := range msgs {
			bytes := d.Body
			if err != nil {
				log.Println(err)
			}
			fmt.Fprintln(file, string(bytes))
		}
	}()
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func getQueue(ch *amqp.Channel) <-chan amqp.Delivery {
	q, err := ch.QueueDeclare(
		"main_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	return msgs
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
