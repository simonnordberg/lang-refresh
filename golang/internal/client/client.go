package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	sessionEndpoint = "/api/v1/session"
)

type Client struct {
	restClient *resty.Client
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponseSuccess struct {
	Id string `json:"id"`
}

type AuthResponseError struct {
	Error string `json:"error"`
}

func NewClient(hostname string) *Client {
	client := resty.New()
	client.SetBaseURL(hostname)
	client.SetDebug(false) // Make configurable
	return &Client{restClient: client}
}

func (c *Client) Authenticate(username, password string) (string, error) {
	payload := AuthRequest{
		Username: username,
		Password: password,
	}

	var responseBody AuthResponseSuccess
	var responseError AuthResponseError

	_, err := c.restClient.R().
		SetBody(&payload).
		SetResult(&responseBody).
		SetError(&responseError).
		Post(sessionEndpoint)

	if err != nil {
		return "", fmt.Errorf("failed to call api: %w", err)
	}

	if responseError.Error != "" {
		return "", fmt.Errorf("error from server: %s", responseError.Error)
	}

	return responseBody.Id, nil
}
