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
	type fields struct {
		repos storage.Repositories
	}

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		fields  fields
		request string
		want    want
	}{
		// TODO: Add test cases.
		{
			name: "Test Content-type and Status Code",
			fields: fields{
				repos: storage.NewLocalStorage(),
			},
			request: "http://127.0.0.1:8080/update/gauge/Alloc/2156",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
		},
		{
			name: "Test Status Code wrong url",
			fields: fields{
				repos: storage.NewLocalStorage(),
			},
			request: "http://127.0.0.1:8080/update/gauger/MAlloc/2156",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
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
			handler := http.HandlerFunc(h.ServeHTTP)

			handler.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, result.StatusCode, tt.want.statusCode)
			assert.Equal(t, result.Header.Get("Content-Type"), tt.want.contentType)

			_, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
