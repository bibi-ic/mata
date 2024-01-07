package helper

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresImage struct {
	Container        *postgres.PostgresContainer
	ConnectionString string
}

func NewPostgresContainer(ctx context.Context) (*PostgresImage, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.1-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "..", "db/migration", "000001_init_schema.up.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresImage{
		Container:        pgContainer,
		ConnectionString: connStr,
	}, nil
}

func (pgI *PostgresImage) Drop(ctx context.Context) {
	if err := pgI.Container.Terminate(ctx); err != nil {
		log.Fatal("error terminating postgres container: ", err)
	}
}
