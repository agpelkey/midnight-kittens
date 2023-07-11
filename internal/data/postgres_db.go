package data

import (
	"context"
	"database/sql"
	"time"
)

type DBRepo interface {
    sendFactToDB(fact *CatFact) error
}


type PostgresDB struct {
    DB *sql.DB
}

func (m PostgresDB) sendFactToDB(fact *CatFact) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `INSERT INTO cat_facts VALUES (fact) VALUES ($1)`

    args := []interface{}{fact.Fact}

    return m.DB.QueryRowContext(ctx, query, args...).Scan(&fact.Fact)
}

