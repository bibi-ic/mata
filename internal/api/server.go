package api

import (
	"context"
	"log"

	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/internal/cache"
	"github.com/bibi-ic/mata/internal/db/seed"
	db "github.com/bibi-ic/mata/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	config config.Config
	router *gin.Engine
	cache  cache.Cache
	store  db.Store
}

func NewServer(cfg config.Config) *Server {

	// Init Database Connection Pool
	connPool, err := pgxpool.New(context.Background(), cfg.DB.Source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	runDBMigration(cfg.MigrationURL, cfg.DB.Source)

	// Insert Bulk Key API
	count, err := seed.Key(connPool, cfg.Iframely.Key)
	if err != nil && db.ErrorCode(err) != db.UniqueViolation {
		log.Fatal("cannot seed key: ", err)
	} else {
		log.Printf("inserted %v keys", count)
	}

	//  Connect to cache server
	opts, err := redis.ParseURL(cfg.Cache.Source)
	if err != nil {
		log.Fatal("cannot parse redis source: ", err)
	}
	rClient := redis.NewClient(opts)
	_, err = rClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("cannot connect to redis: ", err)
	}
	defer rClient.Close()

	cache := cache.New(rClient, cfg.Cache.Age)
	store := db.NewStore(connPool)

	s := &Server{
		config: cfg,
		cache:  cache,
		store:  store,
	}
	s.setUpRoutes()
	return s
}

func (s *Server) setUpRoutes() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	s.router.POST("/meta", s.Retrieve)

	s.router = router
}

func (s *Server) Start() error {
	address := s.config.Server.Address + ":" + s.config.Server.Port
	return s.router.Run(address)
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
