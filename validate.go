package httpc

import (
	"errors"
	"net/url"
)

// validRequest validates the request
func (h *HttpClient) validRequest() (bool, error) {
	if h.Error != nil {
		return false, h.Error
	}
	if h.Method == "" {
		return false, errors.New("no HTTP method specified")
	}
	return true, nil
}

func (h *HttpClient) validURL() (bool, error) {
	_, err := url.ParseRequestURI(h.Url)
	if err != nil {
		return false, err
	}
	return true, nil
}
