package api

import (
	"io"
	"net/http"
)

// NewIframelyRequest generate an API call to 3rd Iframely for fetching meta from link given
func NewIframelyRequest(url, key string) (*http.Request, error) {
	scheme := "https://"
	u := scheme + "iframe.ly/api/oembed"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("url", url)
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// IframelyResponse read response after call API to Iframely for handling data
func IframelyResponse(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return b, nil
}
