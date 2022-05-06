package httpc

import (
	"testing"
)

func TestHttpClient_validURL(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid URL returns true",
			args:    args{endpoint: "https://api.com/api/v1/example/"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid URL returns false",
			args:    args{endpoint: "api.com/api/v1/example"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := Get(tt.args.endpoint)
			got := hc.validURL()
			if got != tt.want {
				t.Errorf("HttpClient.validURL() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && hc.Error == nil {
				t.Error("want URL error set on hc.Error")
			}
		})
	}
}
