package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/bibi-ic/mata/config"
	"github.com/bibi-ic/mata/internal/cache"
	"github.com/bibi-ic/mata/internal/db/seed"
	db "github.com/bibi-ic/mata/internal/db/sqlc"
	"github.com/bibi-ic/mata/internal/opentelemetry"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	router.Use(otelgin.Middleware("mata-server"))

	router.POST("/meta", s.Retrieve)

	s.router = router
}

func (s *Server) Start() (err error) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := opentelemetry.SetupOtelSDK(ctx, s.config.Jaeger.Source)
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Set up HTTP server.
	address := s.config.Server.Address + ":" + s.config.Server.Port
	srv := &http.Server{
		Addr:         address,
		Handler:      s.router,
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
			return
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return
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
