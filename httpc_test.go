package httpc

import (
	"bytes"
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

func TestGetJsonWithDefaultHeaders(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	defaultHeaders := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	hc := GetJson(url)
	if got := hc.headers; !reflect.DeepEqual(got, defaultHeaders) {
		t.Errorf("Get() headers got %v, want %v", got, defaultHeaders)
	}
}

func TestPost(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	type requestBody struct {
		Text  string
		Token string
	}
	body := &requestBody{
		Text:  "this is some text",
		Token: "mySecretToken",
	}
	hc := Post(url, body)
	if got := hc.Method; !reflect.DeepEqual(got, http.MethodPost) {
		t.Errorf("Post() Method got %v, want %v", got, http.MethodPost)
	}
	if got := hc.Url; !reflect.DeepEqual(got, url) {
		t.Errorf("Post() url got %v, want %v", got, url)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(hc.body)
	got := buf.String()
	want := "{\"Text\":\"this is some text\",\"Token\":\"mySecretToken\"}"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Post() request body got %v, want %v", got, want)
	}
}

func TestPostErrorSet(t *testing.T) {
	t.Parallel()
	badUrl := "api.com/api/v1/example/"
	type requestBody struct {
		Text  string
		Token string
	}
	body := &requestBody{
		Text:  "this is some text",
		Token: "mySecretToken",
	}

	hc := Post(badUrl, body)
	if hc.Error == nil {
		t.Error("want error set on hc.Error")
	}
	if hc.Method != "" {
		t.Error("want error set on hc.Error")
	}
	if hc.body != nil {
		t.Error("want error set on hc.Error")
	}
}

func TestAddHeaders(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	headers := map[string]string{
		"X-Auth-Token": "mySecretToken",
	}
	hc := Get(url).AddHeaders(headers)
	if got := hc.headers; !reflect.DeepEqual(got, headers) {
		t.Errorf("Get() headers got %v, want %v", got, headers)
	}
}

func TestAddHeadersErrorSet(t *testing.T) {
	t.Parallel()
	badUrl := "api.com/api/v1/example/"
	headers := map[string]string{
		"X-Auth-Token": "mySecretToken",
	}
	wantHeaders := map[string]string{}
	hc := Get(badUrl).AddHeaders(headers)
	if hc.Error == nil {
		t.Error("want URL error set on hc.Error")
	}
	if got := hc.headers; !reflect.DeepEqual(got, wantHeaders) {
		t.Errorf("AddHeaders() got %v, want %v", got, wantHeaders)
	}
}

func TestBasicAuth(t *testing.T) {
	t.Parallel()
	url := "https://api.com/api/v1/example/"
	username := "user"
	password := "myPassword"
	hc := Get(url).BasicAuth(username, password)
	if got := hc.basicAuthRequired; got != true {
		t.Errorf("BasicAuth() got %v, want true", got)
	}
	if got := hc.basicAuthUsername; got != username {
		t.Errorf("BasicAuth() got %v, want %v", got, username)
	}
	if got := hc.basicAuthPassword; got != password {
		t.Errorf("BasicAuth() got %v, want %v", got, password)
	}
}

func TestBasicAuthErr(t *testing.T) {
	t.Parallel()
	badUrl := "api.com/api/v1/example/"
	username := "user"
	password := "myPassword"
	hc := Get(badUrl).BasicAuth(username, password)
	if got := hc.basicAuthRequired; got != false {
		t.Errorf("BasicAuth() got %v, want false", got)
	}
	if got := hc.basicAuthUsername; got != "" {
		t.Errorf("BasicAuth() got %v, want %v", got, "")
	}
	if got := hc.basicAuthPassword; got != "" {
		t.Errorf("BasicAuth() got %v, want %v", got, "")
	}
}

func TestDo_StatusOK(t *testing.T) {
	t.Parallel()
	tc, mux, teardown := testClient(t)
	t.Cleanup(teardown)
	endpoint := "/api/v1/example/"
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(loadTestJson("testdata/simple.json")))
	})
	username := "user"
	password := "myPassword"
	got, err := Get(tc.Url+endpoint).BasicAuth(username, password).Do()
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
	t.Cleanup(teardown)
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
	t.Cleanup(teardown)
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
