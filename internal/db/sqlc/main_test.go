package db

import (
	"context"
	"log"
	"os"
	"testing"

	testhelpers "github.com/bibi-ic/mata/test/helper"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	// Setup Suite
	ctx := context.Background()
	pgSuit, err := testhelpers.NewPostgresContainer(ctx)
	if err != nil {
		log.Fatal("error cannot create postgres container: ", err)
	}

	connPool, err := pgxpool.New(ctx, pgSuit.ConnectionString)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(connPool)

	code := m.Run()

	// Teardown Suite
	pgSuit.Drop(ctx)

	os.Exit(code)
}
