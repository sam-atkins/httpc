package httpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type HttpClient struct {
	basicAuthRequired bool
	basicAuthUsername string
	basicAuthPassword string
	body              io.Reader
	Client            *http.Client
	Error             error
	headers           map[string]string
	Method            string
	Url               string
}

func NewClient(url string) *HttpClient {
	return &HttpClient{
		body:    nil,
		Client:  &http.Client{},
		headers: map[string]string{},
		Url:     url,
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

// GetJson is a convenience wrapper on the Get method. It sets the Content Type header to JSON.
func GetJson(url string) *HttpClient {
	h := Get(url)
	h.headers = map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	return h
}

// Post prepares a Post JSON request. It sets the http method to Post, validates the provided
// url and sets the content type to JSON. It sets the arg requestBody as the post request
// body.
func Post(url string, requestBody interface{}) *HttpClient {
	h := NewClient(url)
	validUrl, err := h.validURL()
	if !validUrl {
		h.Error = err
		return h
	}
	h.Method = http.MethodPost
	h.headers = map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	body, err := json.Marshal(&requestBody)
	if err != nil {
		h.Error = err
		return h
	}
	h.body = bytes.NewBuffer(body)
	return h
}

type FormField struct {
	Name  string
	Value string
}

type FormFile struct {
	FileName string
	File     *os.File
}

// PostForm prepares a POST form request.
func PostForm(url string, formFile FormFile, formFields ...FormField) *HttpClient {
	h := NewClient(url)
	validUrl, err := h.validURL()
	if !validUrl {
		h.Error = err
		return h
	}

	h.Method = http.MethodPost

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, field := range formFields {
		var formField io.Writer
		formField, errForm := writer.CreateFormField(field.Name)
		if errForm != nil {
			h.Error = errForm
			return h
		}
		_, errCopy := io.Copy(formField, strings.NewReader(field.Value))
		if errCopy != nil {
			h.Error = errCopy
			return h
		}
	}

	fileField, err := writer.CreateFormFile("file", formFile.FileName)
	if err != nil {
		h.Error = err
		return h
	}
	_, err = io.Copy(fileField, formFile.File)
	if err != nil {
		h.Error = err
		return h
	}

	writer.Close()
	h.body = bytes.NewReader(body.Bytes())
	h.headers = map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	return h
}

// AddHeaders adds headers to the request.
func (h *HttpClient) AddHeaders(headers map[string]string) *HttpClient {
	if h.Error != nil {
		return h
	}
	for key, value := range headers {
		h.headers[key] = value
	}
	return h
}

// BasicAuth sets basic auth on the request.
func (h *HttpClient) BasicAuth(username, password string) *HttpClient {
	if h.Error != nil {
		return h
	}
	h.basicAuthRequired = true
	h.basicAuthUsername = username
	h.basicAuthPassword = password
	return h
}

// Do validates the request and if ok, does the HTTP request. It returns the HTTP response
// or an error
func (h *HttpClient) Do() (*http.Response, error) {
	validReq, reqErr := h.validRequest()
	if !validReq {
		return nil, reqErr
	}

	req, err := http.NewRequest(h.Method, h.Url, h.body)
	if err != nil {
		return nil, err
	}

	for key, value := range h.headers {
		req.Header.Add(key, value)
	}

	if h.basicAuthRequired {
		req.SetBasicAuth(h.basicAuthUsername, h.basicAuthPassword)
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

// Load makes the HTTP request and unmarshals the response into the provided data arg.
// This method calls Do() so there is no need to call Do then Load.
func (h *HttpClient) Load(data interface{}) error {
	res, err := h.Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	decodeErr := json.NewDecoder(res.Body).Decode(&data)
	if decodeErr != nil {
		return decodeErr
	}
	return nil
}
