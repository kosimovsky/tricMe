package server

import (
	"net/http"
	"testing"
)

func TestServer_Run(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	type fields struct {
		httpServer *http.Server
	}
	type args struct {
		port    string
		handler func(w http.ResponseWriter, response *http.Request)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test status code",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				httpServer: tt.fields.httpServer,
			}
			if err := s.Run(tt.args.port, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
