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
	"net/url"
	"time"
)

// TODO Added userRequestID into web api
type WebClient struct {
	baseURL string
	active  bool
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

func (c WebClient) IsActive() bool {
	return c.active
}

func (c WebClient) SearchAudio(ctx context.Context, option voice.AudioSearch, page uint) ([]api.AudioFileInfo, error) {
	searchType := "id"
	if option.Type == voice.NameType {
		searchType = "name"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(
		"%s/api/v1/audio/fileinfo/%s/%s?page=%d",
		c.baseURL,
		searchType,
		option.Value,
		page,
	), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err == nil || res.StatusCode != 200 {
		err = fmt.Errorf("failed download audio. Response: %s", body)
	}

	var result api.GetAudioIdsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Files, nil
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

func NewClient(baseURL string) (WebClient, error) {
	if _, err := url.Parse(baseURL); err != nil {
		return WebClient{}, err
	}

	client := WebClient{
		baseURL: baseURL,
		client: http.Client{
			Timeout: time.Second * 5,
		},
	}

	slog.Info("Testing connection with Web API")
	if err := client.ping(); err != nil {
		return WebClient{}, err
	}

	client.active = true

	return client, nil
}
