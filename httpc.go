package httpc

import (
	"errors"
	"io"
	"net/http"
)

type HttpClient struct {
	Body    io.Reader
	Client  *http.Client
	Error   error
	headers map[string]string
	Method  string
	Url     string
}

func NewClient(url string) *HttpClient {
	return &HttpClient{
		Body:   nil,
		Client: &http.Client{},
		headers: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		Url: url,
	}
}

// Get prepares a GET request. It sets the http method to GET and validates the provided
// url
func Get(url string) *HttpClient {
	h := NewClient(url)
	validUrl, err := h.validURL()
	if !validUrl {
		h.Error = err
		return h
	}
	h.Method = http.MethodGet
	return h
}

// AddHeaders adds headers to the request.
//
// Defaults already set are:
//   Content-Type: application/json;charset=UTF-8
func (h *HttpClient) AddHeaders(headers map[string]string) *HttpClient {
	if h.Error != nil {
		return h
	}
	for key, value := range headers {
		h.headers[key] = value
	}
	return h
}

// Do validates the request and if ok, does the HTTP request. It returns the HTTP response
// or an error
func (h *HttpClient) Do() (*http.Response, error) {
	validReq, reqErr := h.validRequest()
	if !validReq {
		return nil, reqErr
	}

	req, err := http.NewRequest(h.Method, h.Url, h.Body)
	if err != nil {
		return nil, err
	}

	for key, value := range h.headers {
		req.Header.Add(key, value)
	}

	res, resErr := h.Client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	// if not 2xx
	if (res.StatusCode != http.StatusOK) &&
		(res.StatusCode != http.StatusCreated) &&
		(res.StatusCode != http.StatusAccepted) &&
		(res.StatusCode != http.StatusNoContent) {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(resBody))
	}
	// 2xx
	return res, nil
}
