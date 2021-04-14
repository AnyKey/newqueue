package main

import (
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

	forever := make(chan bool)
	f, err := os.OpenFile("main.log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	go func() {
		for d := range msgs {
			//var res model.LogBody
			//err = json.Unmarshal(d.Body, &res)
			//if err != nil {
			//	log.Println(err)
			//}
			//newString  := fmt.Sprintf("name: %v | action: %v | time: %v \n", res.Name, res.Action, res.Time)
			bytes := d.Body
			_, err = f.Write(bytes)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
