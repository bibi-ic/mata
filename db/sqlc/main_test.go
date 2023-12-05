package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/bibi-ic/mata/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	c, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), c.DB.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
