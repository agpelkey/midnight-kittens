package main

import (
	"flag"
	"fmt"
	"log"

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

	// declare application instance
	app := &application{
		config: cfg,
	}

	fmt.Println(app.handleGetCatFact())

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}

}
