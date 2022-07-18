package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosimovsky/tricMe/internal/storage"
)

func TestHandler_MetricsRouterRouter(t *testing.T) {
	type handler struct {
		repos storage.Repositories
	}
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		fields     handler
		path       string
		method     string
		statusCode int
	}{
		// TODO: Add test cases.
		{
			name: "Test good Status Code",
			fields: handler{
				repos: storage.NewLocalStorage().DefaultStorage(),
			},
			statusCode: http.StatusOK,
			path:       "/update/gauge/Alloc/98479",
			method:     "POST",
		},
		{
			name: "Test without value",
			fields: handler{
				repos: storage.NewLocalStorage().DefaultStorage(),
			},
			statusCode: http.StatusNotFound,
			path:       "/update/gauge/Alloc",
			method:     "POST",
		},
		{
			name: "get gauge value",
			fields: handler{
				repos: storage.NewLocalStorage().DefaultStorage(),
			},
			statusCode: http.StatusOK,
			path:       "/value/gauge/Alloc",
			method:     "GET",
		},
		{
			name: "get start page",
			fields: handler{
				repos: storage.NewLocalStorage().DefaultStorage(),
			},
			statusCode: http.StatusOK,
			path:       "/",
			method:     "GET",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				repos: tt.fields.repos,
			}

			ts := httptest.NewServer(h.MetricsRouter())
			defer ts.Close()

			resp, _ := testRequest(t, ts, tt.method, tt.path)
			assert.Equal(t, resp.StatusCode, tt.statusCode)
			resp.Body.Close()
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
