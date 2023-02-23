package httpHelper

import (
	"io"
	"net/http"
	"net/url"
)

func Get(url string, header http.Header, query url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	return http.DefaultClient.Do(req)
}

func Post(url string, body io.Reader, header http.Header, query url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	return http.DefaultClient.Do(req)
}
