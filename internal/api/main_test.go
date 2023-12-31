package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func makeGinContext(recorder *httptest.ResponseRecorder) *gin.Context {
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx
}

func mockJsonRetrieve(c *gin.Context, u url.Values) {
	c.Request.Method = http.MethodPost
	c.Request.Header.Set("Content-Type", "application/json")

	c.Request.URL.RawQuery = u.Encode()
}
