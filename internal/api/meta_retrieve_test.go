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

	mockcache "github.com/bibi-ic/mata/internal/cache/mock"
	mockdb "github.com/bibi-ic/mata/internal/db/mock"
	"github.com/bibi-ic/mata/internal/models"
	"github.com/bibi-ic/mata/internal/utils"
	"github.com/jarcoal/httpmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TODO: fix test API layer next merge

func TestRetrieveAPI(t *testing.T) {
	key := utils.RandKey()
	u := utils.RandURL()

	meta, err := randomMeta(u)
	require.NoError(t, err)

	scheme := "https://"
	urlRequest := scheme + "iframe.ly/api/oembed"

	testCases := []struct {
		name             string
		buildStubs       func(store *mockdb.MockStore, cache *mockcache.MockCache)
		fakeExternalCall func()
		checkResult      func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			name: "201 Created MetaInserted",
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {

				store.EXPECT().
					GetAPITx(gomock.Any()).
					Times(1).
					Return(key, nil)

				cache.EXPECT().
					Get(gomock.Any(), u).
					Times(1).
					Return(nil, redis.Nil)

				cache.EXPECT().
					Set(gomock.Any(), u, &meta).
					Times(1).
					Return(nil)
			},
			fakeExternalCall: func() {
				baseURL, err := url.Parse(urlRequest)
				require.NoError(t, err)

				params := url.Values{}
				params.Add("url", u)
				params.Add("key", key)

				baseURL.RawQuery = params.Encode()

				httpmock.RegisterResponder("GET", baseURL.String(),
					func(req *http.Request) (*http.Response, error) {
						return httpmock.NewJsonResponse(http.StatusOK, meta)
					},
				)
			},
			checkResult: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, r.Code)
				requireBodyMatchMeta(t, r.Body, meta)
			},
		},
		{
			name: "200 OK MetaFound InCache",
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache) {
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
			fakeExternalCall: func() {},
			checkResult: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
				requireBodyMatchMeta(t, r.Body, meta)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			store := mockdb.NewMockStore(ctrl)
			cache := mockcache.NewMockCache(ctrl)
			tc.buildStubs(store, cache)
			tc.fakeExternalCall()

			recorder := httptest.NewRecorder()
			ctx := makeGinContext(recorder)

			q := url.Values{}
			q.Add("url", u)

			mockJsonRetrieve(ctx, q)
			server := newTestServer(t, store, cache)
			server.Retrieve(ctx)

			tc.checkResult(t, recorder)
		})
	}
}

func randomMeta(url string) (models.Meta, error) {
	var gotMeta = utils.RandMeta(url)

	err := gotMeta.Parse()
	return gotMeta, err
}

func randomIframelyKey(keys []string) string {
	return keys[rand.Intn(len(keys))]
}

func requireBodyMatchMeta(t *testing.T, body *bytes.Buffer, meta models.Meta) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMeta models.Meta
	err = json.Unmarshal(data, &gotMeta)
	require.NoError(t, err)
	require.Equal(t, meta, gotMeta)
}
