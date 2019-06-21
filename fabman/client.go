package fabman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://fabman.io/api/v1"

// Client abstracts FabMan APIs
type Client struct {
	account    int
	token      string
	httpClient *http.Client
}

// NewClient returns new FabMan API client
func NewClient(account int, token string) *Client {
	return &Client{
		account:    account,
		token:      token,
		httpClient: http.DefaultClient,
	}
}

func (client *Client) request(method, path string, src, dst interface{}) error {
	var body io.Reader

	if src != nil {
		payload, err := json.Marshal(src)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(payload)
	}

	url := fmt.Sprintf("%s/%s", baseURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+client.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var e Error
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}
		return e
	}

	if dst != nil {
		return json.NewDecoder(resp.Body).Decode(dst)
	}

	return nil
}

func (client *Client) get(path string, dst interface{}) error {
	return client.request(http.MethodGet, path, nil, dst)
}

func (client *Client) create(path string, obj interface{}) error {
	return client.request(http.MethodPost, path, obj, obj)
}

func (client *Client) update(path string, obj interface{}) error {
	return client.request(http.MethodPut, path, obj, obj)
}

func (client *Client) delete(path string) error {
	return client.request(http.MethodDelete, path, nil, nil)
}
