package main

import (
	"context"
	"log"

	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/db/seed"
	db "github.com/bibi-ic/mata/db/sqlc"
	"github.com/bibi-ic/mata/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), c.DB.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// Insert Bulk Key API
	count, err := seed.Key(connPool, c.Iframely.Key)
	if err != nil && db.ErrorCode(err) != db.UniqueViolation {
		log.Fatal("cannot seed key: ", err)
	} else {
		log.Printf("inserted %v keys", count)
	}

	store := db.NewStore(connPool)
	s := server.New(c, store)
	err = s.Start()
	if err != nil {
		log.Fatal("can not run server: ", err)
	}
}
