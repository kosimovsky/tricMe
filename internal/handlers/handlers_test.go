package handlers

import (
	"github.com/kosimovsky/tricMe/internal/storage"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	type handler struct {
		repos storage.Repositories
	}

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		fields  handler
		request string
		want    want
	}{
		// TODO: Add test cases.
		{
			name: "Test good Status Code",
			fields: handler{
				repos: storage.NewLocalStorage(),
			},
			request: "http://127.0.0.1:8080/update/gauge/Alloc/2156",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
		},
		{
			name: "Test Status Code wrong type of metrics",
			fields: handler{
				repos: storage.NewLocalStorage(),
			},
			request: "http://127.0.0.1:8080/update/gauger/Alloc/2156",
			want: want{
				contentType: "text/plain",
				statusCode:  501,
			},
		},
		{
			name: "Test Status Code: without metric",
			fields: handler{
				repos: storage.NewLocalStorage(),
			},
			request: "http://127.0.0.1:8080/update/counter/",
			want: want{
				contentType: "text/plain",
				statusCode:  404,
			},
		},
		{
			name: "Test Status Code: any metric with none",
			fields: handler{
				repos: storage.NewLocalStorage(),
			},
			request: "http://127.0.0.1:8080/update/gauge/AnyCounter/none",
			want: want{
				contentType: "text/plain",
				statusCode:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			h := Handler{
				repos: tt.fields.repos,
			}

			hdlr := http.HandlerFunc(h.MetricsHandler)

			hdlr.ServeHTTP(w, request)

			result := w.Result()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
