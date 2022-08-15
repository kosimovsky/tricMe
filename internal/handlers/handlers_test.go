package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosimovsky/tricMe/internal/storage"
)

func TestHandler_MetricsRouterRouter(t *testing.T) {

	store := storage.TestMetrics()
	_ = store.Restore("../test/test.json", true)
	handler := NewHandler(store)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		path        string
		method      string
		body        string
		contentType string
		statusCode  int
	}{
		// TODO: Add test cases.
		{
			name:        "Test good Status Code",
			statusCode:  http.StatusOK,
			path:        "/value/",
			contentType: "application/json",
			body:        `{"id": "Alloc","type": "gauge"}`,
			method:      http.MethodPost,
		},
		{
			name:        "Test incorrect type",
			statusCode:  http.StatusNotFound,
			path:        "/value/",
			method:      http.MethodPost,
			body:        `{"id": "RandomValue","type": "gauger"}`,
			contentType: "application/json",
		},
		{
			name:       "get gauge value",
			statusCode: http.StatusOK,
			path:       "/value/gauge/Alloc",
			method:     "GET",
		},
		{
			name:       "get start page",
			statusCode: http.StatusOK,
			path:       "/",
			method:     "GET",
		},
	}

	ts := httptest.NewServer(handler.MetricsRouter())
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.method, tt.path, tt.body, tt.contentType)
			assert.Equal(t, tt.statusCode, resp.StatusCode)
			resp.Body.Close()
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, body, contentType string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	req.Header.Set("Content-Type", contentType)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
