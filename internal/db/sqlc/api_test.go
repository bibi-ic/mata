package db

import (
	"context"
	"testing"

	"github.com/bibi-ic/mata/internal/randomutil"
	"github.com/stretchr/testify/require"
)

func createRandomAPI(t *testing.T) Api {
	arg := randomutil.RandKey()

	api, err := testStore.CreateAPI(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, api)

	require.Equal(t, arg, api.Key)
	require.Zero(t, api.UsageCount)

	require.NotZero(t, api.ID)
	require.NotZero(t, api.CreatedAt)
	return api
}

func TestCreateAPI(t *testing.T) {
	createRandomAPI(t)
}

func TestGetAPI(t *testing.T) {
	createRandomAPI(t)
	api, err := testStore.GetAPI(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, api)

	require.NotEmpty(t, api.Key)
	require.GreaterOrEqual(t, api.UsageCount, int64(0))

	require.NotZero(t, api.ID)
	require.NotZero(t, api.CreatedAt)
}
