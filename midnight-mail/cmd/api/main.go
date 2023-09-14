package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	//"net/smtp"
	//"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rabbitmq/amqp091-go"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
        os.Exit(1)
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

    for payload := range msgs {
        log.Printf("Received fact: %s", payload.Body)

        reader := strings.NewReader(string(payload.Body))

        var bucket = "midnight-cat-fact"
        var accessKey = os.Getenv("AWS_ACCESS_KEY")
        //var secretKey = os.Getenv("AWS_SECRET_KEY")
        //var filename = os.File("~/.aws/config")

        sess, err := session.NewSession(&aws.Config{
            Region: aws.String("us-west-1"),
        })

        uploader := s3manager.NewUploader(sess)

        _, err = uploader.Upload(&s3manager.UploadInput{
            Bucket: aws.String(bucket),
            Key: &accessKey,
            Body: reader,
        })
        if err != nil {
            log.Fatal("%v", err)
        }
    }

    fmt.Println("successfully uploaded to bucket")

    /*
    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
        }
    }()

    //log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    */

    <-forever

}





