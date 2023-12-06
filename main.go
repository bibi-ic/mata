package main

import (
	"context"
	"log"

	"github.com/bibi-ic/mata/cache"
	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/db/seed"
	db "github.com/bibi-ic/mata/db/sqlc"
	"github.com/bibi-ic/mata/server"
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
	s := server.New(c, store, cache)
	err = s.Start()
	if err != nil {
		log.Fatal("can not run server: ", err)
	}
}
