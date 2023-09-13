package main

import (
	"log"
	"os"
	//"net/smtp"
	//"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

// function to populate .env variables
func LoadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
        os.Exit(1)
    }
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
        os.Exit(1)
	}
}

func main() {

    LoadEnv()

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


	//var forever chan struct{}

    for payload := range msgs {
        log.Printf("Received fact: %s", payload.Body)
    }

    /*
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    */

	//<-forever
}


/*
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
*/



