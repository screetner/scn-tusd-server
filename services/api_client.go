package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
)

var (
	clientInstance *APIClient
	once           sync.Once
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func GetAPIClient() *APIClient {
	once.Do(func() {
		backendURL := os.Getenv("API_URL")
		if backendURL == "" {
			println("no service account key for Screetner backend user using the API_URl environment variable")
		}

		clientInstance = &APIClient{
			BaseURL:    backendURL,
			HTTPClient: &http.Client{},
		}
	})
	return clientInstance
}

func (client *APIClient) Get(endpoint string, headers map[string]string, body io.Reader) (*http.Response, error) {
	url := client.BaseURL + endpoint
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *APIClient) Post(endpoint string, data interface{}) (*http.Response, error) {
	url := client.BaseURL + endpoint
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := client.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
