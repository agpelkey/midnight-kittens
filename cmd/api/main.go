package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/agpelkey/midnight-kittens/internal/data"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
)

// declare config
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// declare application struct
type application struct {
	config config
	facts  data.CatFact
    Models data.Models
}

// runnit
func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API Server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

    flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("CAT_FACT_DB_DSN"), "PostgreSQL DB DSN")

	// Pass in values to optimize our database
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connections idle time")

	flag.Parse()

	// connect to db here
    conn, err := openDB(cfg)
    if err != nil {
        log.Fatal(err)
    }

    defer conn.Close()

	// declare application instance
	app := &application{
		config: cfg,
        Models: data.NewModels(conn), 
	}

    //fmt.Println(app.handleGetCatFact())
    payload, err := app.handleGetCatFact()
    if err != nil {
        log.Fatal(err)
    }

    createRabbitPublisher(payload)

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
    
    // use sql.open to create connection pool
    db, err := sql.Open("postgres", cfg.db.dsn)
    if err != nil {
        return nil, err
    }

    // The next chunk of code is to establish our database connection config.
    db.SetMaxOpenConns(cfg.db.maxOpenConns)

    db.SetMaxIdleConns(cfg.db.maxIdleConns)

    // For the last one, we need to use time.ParseDuration() to convert idle timeout
    // duration string to a time.Duration type
    duration, err := time.ParseDuration(cfg.db.maxIdleTime)
    if err != nil {
        return nil, err 
    }

    // now set max idle time
    db.SetConnMaxIdleTime(duration)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = db.PingContext(ctx)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}

func createRabbitPublisher(fact *data.CatFact) {
    
    conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "cat_fact",
        false,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "failed to declare a queue")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body := fact.Fact

    err = ch.PublishWithContext(ctx,
        "", // exchange
        q.Name, // routing key
        false, // mandatory
        false, // immediate
        amqp091.Publishing{
            ContentType: "text/plain",
            Body: []byte(body),
        })

    failOnError(err, "failed to publish a message")
    log.Printf(" [x] Sent %s\n", body)

}









