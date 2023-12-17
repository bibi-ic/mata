package seed

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Key(connPool *pgxpool.Pool, k []string) (int64, error) {
	rows := make([][]interface{}, 0)
	for _, v := range k {
		rows = append(rows, []interface{}{v})
	}

	copyCount, err := connPool.CopyFrom(
		context.Background(),
		pgx.Identifier{"apis"},
		[]string{"key"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return 0, err
	}
	return copyCount, nil
}
