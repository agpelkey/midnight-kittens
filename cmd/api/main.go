package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	"github.com/agpelkey/midnight-kittens/internal/data"
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
