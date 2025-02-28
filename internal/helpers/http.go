package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func DefineHeaders(apiKey string) http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return headers
}

func CreateQueryString(params map[string]string) string {
	query := url.Values{}
	for key, value := range params {
		query.Set(key, value)
	}
	return query.Encode()
}

func Request(url string, method string, headers http.Header, body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
