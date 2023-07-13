package data

import (
	"context"
	"database/sql"
	"time"
)

type DBRepo interface {
    SendFactToDB(fact *CatFact) error
}


type PostgresDB struct {
    DB *sql.DB
}

func (m PostgresDB) SendFactToDB(fact *CatFact) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `INSERT INTO cat_facts (fact, length) VALUES ($1, $2)`

    args := []interface{}{fact.Fact}

    err := m.DB.QueryRowContext(ctx, query, args...).Scan(fact.Fact, fact.Length)
    if err != nil {
        return err
    }

    return nil 
}

