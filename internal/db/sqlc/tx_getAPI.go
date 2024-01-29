package db

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetAPITx performs an API using
// It randomly selects an API, and update usage count within a database transaction
func (store *SQLStore) GetAPITx(ctx context.Context) (string, error) {
	var result string

	ctx, span := trace.SpanFromContext(ctx).TracerProvider().
		Tracer("internal/db/sqlc").
		Start(ctx, "getAPITx")

	defer span.End()
	span.SetAttributes(attribute.String("db.system", "postgresql"))

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var a Api
		a, err = q.GetAPI(ctx)
		if err != nil {
			return err
		}

		err = q.UpdateAPIUsageCount(ctx, UpdateAPIUsageCountParams{
			ID:         a.ID,
			UsageCount: a.UsageCount + 1,
		})
		if err != nil {
			return err
		}

		result = a.Key
		return err
	})
	return result, err
}
