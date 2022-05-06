package httpc

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Test Fixtures
func testClient(t *testing.T) (*HttpClient, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(server.URL)

	return client, mux, func() {
		server.Close()
	}
}

func loadTestJson(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	return content
}

// Tests
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

func TestDo_StatusOK(t *testing.T) {
	t.Parallel()
	tc, mux, teardown := testClient(t)
	defer teardown()
	endpoint := "/api/v1/example/"
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(loadTestJson("testdata/simple.json")))
	})

	got, err := Get(tc.Url + endpoint).Do()
	if err != nil {
		t.Errorf("Do() error = %v, wantErr nil", err)
	}
	if got.StatusCode != http.StatusOK {
		t.Errorf("Do() HTTP Status Code = %v, wantErr %v", got.StatusCode, http.StatusOK)
	}
	if got.Body == nil {
		t.Error("Do() want Response.Body, got  nil")
	}
}

func TestDo_StatusNotOK(t *testing.T) {
	t.Parallel()
	tc, mux, teardown := testClient(t)
	defer teardown()
	endpoint := "/api/v1/example/"
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(loadTestJson("testdata/404.json")))
	})

	_, err := Get(tc.Url + endpoint).Do()
	if err == nil {
		t.Error("Do() want error when not 2xx response")
	}
}

func TestDo_InvalidRequestBadURL(t *testing.T) {
	badUrl := "api/v1/example"
	_, err := Get(badUrl).Do()
	if err == nil {
		t.Error("Do() want error when invalid request")
	}
}

func TestDo_InvalidRequestNoMethod(t *testing.T) {
	url := "https://api.com/api/v1/example/"
	hc := Get(url)
	hc.Method = ""
	_, err := hc.Do()
	if err == nil {
		t.Error("Do() want error when invalid request")
	}
}

func TestLoad_StatusOK(t *testing.T) {
	t.Parallel()
	tc, mux, teardown := testClient(t)
	defer teardown()
	endpoint := "/api/v1/example/"
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(loadTestJson("testdata/simple.json")))
	})

	type simpleJSON struct {
		Data []struct {
			ExampleKey string `json:"exampleKey"`
		} `json:"data"`
		Status string `json:"status"`
	}
	var sj simpleJSON

	err := Get(tc.Url + endpoint).Load(&sj)
	if err != nil {
		t.Errorf("Do() error = %v, wantErr nil", err)
	}
	if sj.Status != "OK" {
		t.Errorf("Load() want \"OK\", got %v", sj.Status)
	}
	if len(sj.Data) != 1 {
		t.Errorf("Load() want Data struct len 1, got %v", len(sj.Data))
	}
}

func TestLoad_Err(t *testing.T) {
	t.Parallel()
	type simpleJSON struct {
		Data []struct {
			ExampleKey string `json:"exampleKey"`
		} `json:"data"`
		Status string `json:"status"`
	}
	var sj simpleJSON
	badUrl := "api/v1/example"

	err := Get(badUrl).Load(&sj)
	if err == nil {
		t.Error("Do() want error when invalid request")
	}
}
