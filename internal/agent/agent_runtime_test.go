package agent

import (
	"reflect"
	"testing"
)

func Test_urlGenerator(t *testing.T) {
	type args struct {
		conf config
		m    map[string]gauge
	}
	tests := []struct {
		name     string
		args     args
		wantUrls []string
	}{
		// TODO: Add test cases.
		{
			name: "Validate urls",
			args: args{
				conf: config{
					server: "127.0.0.1",
					port:   "8080",
				},
				m: map[string]gauge{
					"Alloc":       2156,
					"BuckHashSys": 0.2665454,
					"PollCount":   0,
				},
			},
			wantUrls: []string{
				"http://127.0.0.1:8080/update/gauge/Alloc/2156",
				"http://127.0.0.1:8080/update/gauge/BuckHashSys/0.2665454",
				"http://127.0.0.1:8080/update/gauge/PollCount/0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUrls := urlGenerator(tt.args.conf, tt.args.m); !reflect.DeepEqual(gotUrls, tt.wantUrls) {
				t.Errorf("urlGenerator() = %v, want %v", gotUrls, tt.wantUrls)
			}
		})
	}
}
