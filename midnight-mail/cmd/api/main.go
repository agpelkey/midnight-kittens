package main

import (
	"log"
<<<<<<< HEAD
=======
	"net/smtp"
	"os"

>>>>>>> 04ae007 (clean up email api)

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cat_fact", //name
		false,      // durable
		false,      // delete when used
		false,      // exlusive
		false,      // no wait
		nil,        //arguments
	)
	failOnError(err, "failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consume
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to regist a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
<<<<<<< HEAD
=======


func SendMail(fact []byte) error {
    username := os.Getenv("MAILTRAP_USRNAME") 

    password := os.Getenv("MAILTRAP_PASSWD")

    smtpHost := "sandbox.smtp.mailtrap.io"

    auth := smtp.PlainAuth("", username, password, smtpHost)

    // Message data 

    from := username

    to := []string{"agpelkey94@gmail.com"}

    message := fact

    smtpUrl := smtpHost + ":25"

    err := smtp.SendMail(smtpUrl, auth, from, to, message)
    if err != nil {
        log.Fatal(err)
    }

    return nil
}




>>>>>>> 04ae007 (clean up email api)
