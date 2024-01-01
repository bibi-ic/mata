package api

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bibi-ic/mata/config"
	mockcache "github.com/bibi-ic/mata/internal/cache/mock"
	mockdb "github.com/bibi-ic/mata/internal/db/mock"
	"github.com/bibi-ic/mata/internal/external"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/bibi-ic/mata/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRetrieveAPI(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	key := randomIframelyKey(cfg.Iframely.Key)
	u := utils.RandURL()

	meta, err := randomMeta(key, u)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		buildStubs  func(store *mockdb.MockStore, cache *mockcache.MockMataCache)
		checkResult func(t *testing.T, r *httptest.ResponseRecorder, gotMeta *models.Meta)
	}{
		{
			name: "201 Created MetaInserted",
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
			checkResult: func(t *testing.T, r *httptest.ResponseRecorder, gotMeta *models.Meta) {
				require.Equal(t, &meta, gotMeta)
				require.Equal(t, http.StatusCreated, r.Code)
			},
		},
		{
			name: "200 OK MetaFound InCache",
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
			checkResult: func(t *testing.T, r *httptest.ResponseRecorder, gotMeta *models.Meta) {
				require.Equal(t, &meta, gotMeta)
				require.Equal(t, http.StatusOK, r.Code)
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
			server := newTestServer(t, store, cache)
			server.Retrieve(ctx)

			tc.checkResult(t, recorder, &meta)
		})
	}
}

func randomMeta(key, url string) (models.Meta, error) {
	var gotMeta models.Meta
	r, err := external.NewIframelyRequest(url, key)
	if err != nil {
		return gotMeta, err
	}

	res, err := external.IframelyResponse(r)
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

func requireBodyMatchMeta(t *testing.T, body *bytes.Buffer, meta *models.Meta) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMeta *models.Meta
	err = json.Unmarshal(data, gotMeta)
	require.NoError(t, err)
	require.Equal(t, meta, gotMeta)
}
