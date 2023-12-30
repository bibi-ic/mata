package controller

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bibi-ic/mata/api"
	"github.com/bibi-ic/mata/config"
	mockcache "github.com/bibi-ic/mata/internal/cache/mock"
	mockdb "github.com/bibi-ic/mata/internal/db/mock"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/bibi-ic/mata/internal/randomutil"
	"github.com/bibi-ic/mata/internal/status"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRetrieveAPI(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	key := randomIframelyKey(cfg.Iframely.Key)
	u := randomutil.RandURL()

	meta, err := randomMeta(key, u)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		buildStubs  func(store *mockdb.MockStore, cache *mockcache.MockMataCache)
		checkResult func(t *testing.T, meta *models.Meta, st status.Status)
	}{
		{
			name: "201 Created Meta Inserted",
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockMataCache) {
				store.EXPECT().
					GetAPITx(gomock.Any()).
					Times(1).
					Return(key, nil)

				cache.EXPECT().
					Get(gomock.Any(), u).
					Times(1).
					Return(nil, redis.Nil)

				cache.EXPECT().
					Set(gomock.Any(), meta.URL, &meta).
					Times(1).
					Return(nil)
			},
			checkResult: func(t *testing.T, m *models.Meta, st status.Status) {
				require.Equal(t, &meta, m)
				require.Equal(t, http.StatusCreated, st.Code)
			},
		},
		{
			name: "200 OK Meta Found In Cache",
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockMataCache) {
				store.EXPECT().
					GetAPITx(gomock.Any()).
					Times(1).
					Return(key, nil)

				cache.EXPECT().
					Get(gomock.Any(), u).
					Times(1).
					Return(&meta, nil)

				cache.EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResult: func(t *testing.T, m *models.Meta, st status.Status) {
				require.Equal(t, &meta, m)
				require.Equal(t, http.StatusOK, st.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			cache := mockcache.NewMockMataCache(ctrl)
			tc.buildStubs(store, cache)

			recorder := httptest.NewRecorder()
			ctx := makeGinContext(recorder)

			q := url.Values{}
			q.Add("url", u)

			mockJsonRetrieve(ctx, q)

			metaController := NewMetaController(cache, store)
			gotMeta, gotStatus := metaController.Retrieve(ctx, u)

			tc.checkResult(t, gotMeta, gotStatus)
		})
	}
}

func randomMeta(key, url string) (models.Meta, error) {
	var gotMeta models.Meta
	r, err := api.NewIframelyRequest(url, key)
	if err != nil {
		return gotMeta, err
	}

	res, err := api.IframelyResponse(r)
	if err != nil {
		return gotMeta, err
	}

	err = json.Unmarshal(res, &gotMeta)
	if err != nil {
		return gotMeta, err
	}

	err = gotMeta.Parse()
	return gotMeta, err
}

func randomIframelyKey(keys []string) string {
	return keys[rand.Intn(len(keys))]
}
