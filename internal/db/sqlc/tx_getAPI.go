package db

import (
	"context"
)

// GetAPITx performs an API using
// It randomly selects an API, and update usage count within a database transaction
func (store *SQLStore) GetAPITx(ctx context.Context) (string, error) {
	var result string

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
