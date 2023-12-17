package main

import (
	"context"
	"log"

	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/internal/app"
	"github.com/bibi-ic/mata/internal/cache"
	"github.com/bibi-ic/mata/internal/db/seed"
	db "github.com/bibi-ic/mata/internal/db/sqlc"
	"github.com/bibi-ic/mata/internal/server"
	"github.com/bibi-ic/mata/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// Init Database Connection Pool
	connPool, err := pgxpool.New(context.Background(), c.DB.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	runDBMigration(c.MigrationURL, c.DB.Source)

	// Insert Bulk Key API
	count, err := seed.Key(connPool, c.Iframely.Key)
	if err != nil && db.ErrorCode(err) != db.UniqueViolation {
		log.Fatal("cannot seed key: ", err)
	} else {
		log.Printf("inserted %v keys", count)
	}

	//  Connect to cache server
	opts, err := redis.ParseURL(c.Cache.Source)
	if err != nil {
		log.Fatal("cannot parse redis source: ", err)
	}
	rClient := redis.NewClient(opts)
	_, err = rClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("cannot connect to redis: ", err)
	}
	defer rClient.Close()

	cache := cache.New(rClient, c.Cache.Age)
	store := db.NewStore(connPool)

	// Register service
	metaService := service.NewMetaService(cache, store)

	// Starting HTTP server
	s := server.New(c)
	s.RegisterService(app.NewService(metaService))

	err = s.Start()
	if err != nil {
		log.Fatal("can not run server: ", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated successfully")
}
