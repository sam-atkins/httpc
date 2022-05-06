package httpc

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	hc := Get(url)
	if got := hc.Method; !reflect.DeepEqual(got, http.MethodGet) {
		t.Errorf("Get() Method got %v, want %v", got, http.MethodGet)
	}
	if got := hc.Url; !reflect.DeepEqual(got, url) {
		t.Errorf("Get() url got %v, want %v", got, url)
	}
}

func TestGetWithDefaultHeaders(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	defaultHeaders := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	hc := Get(url)
	if got := hc.headers; !reflect.DeepEqual(got, defaultHeaders) {
		t.Errorf("Get() headers got %v, want %v", got, defaultHeaders)
	}
}

func TestAddHeaders(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	headers := map[string]string{
		"X-Auth-Token": "mySecretToken",
	}
	wantHeaders := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
		"X-Auth-Token": "mySecretToken",
	}
	hc := Get(url).AddHeaders(headers)
	if got := hc.headers; !reflect.DeepEqual(got, wantHeaders) {
		t.Errorf("Get() headers got %v, want %v", got, wantHeaders)
	}
}

func TestAddHeadersErrorSet(t *testing.T) {
	t.Parallel()
	badUrl := "api.com/api/v1/example/"
	headers := map[string]string{
		"X-Auth-Token": "mySecretToken",
	}
	wantHeaders := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	hc := Get(badUrl).AddHeaders(headers)
	if hc.Error == nil {
		t.Error("want URL error set on hc.Error")
	}
	if got := hc.headers; !reflect.DeepEqual(got, wantHeaders) {
		t.Errorf("AddHeaders() got %v, want %v", got, wantHeaders)
	}
}
