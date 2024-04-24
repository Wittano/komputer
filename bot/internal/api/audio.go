package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wittano/komputer/api"
	"github.com/wittano/komputer/bot/internal/voice"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// TODO Added userRequestID into web api
type WebClient struct {
	baseURL string
	client  http.Client
	cache   voice.BotLocalStorage
}

func (c WebClient) DownloadAudio(id string) (path string, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/audio/%s", c.baseURL, id), nil)
	if err != nil {
		return "", err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err == nil {
			err = fmt.Errorf("failed download audio. Response: %s", body)
		}

		return "", err
	}

	return c.cache.Add(context.Background(), res.Body, id, res.Header.Get("filename"))
}

func (c WebClient) ping() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/health", c.baseURL), nil)
	if err != nil {
		return err
	}

	c.client.Timeout = time.Millisecond * 500

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("web api isn;t healthy")
	}

	return nil
}

func (c WebClient) Close() error {
	c.client.CloseIdleConnections()

	return nil
}

func NewClient(url string) *WebClient {
	client := WebClient{
		baseURL: url,
		client: http.Client{
			Timeout: time.Second * 5,
		},
	}

	slog.Info("Testing connection with Web API")
	if err := client.ping(); err != nil {
		return nil
	}

	return &client
}
