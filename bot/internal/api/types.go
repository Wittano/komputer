package api

import (
	"fmt"
	"github.com/wittano/komputer/bot/internal/voice"
	"io"
	"net/http"
	"os"
	"time"
)

type WebClient struct {
	baseURL string
	client  http.Client
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

	dest := voice.Path(id)
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err = io.Copy(f, res.Body); err != nil {
		return "", err
	}

	return dest, nil
}

func NewClient(url string) *WebClient {
	return &WebClient{
		baseURL: url,
		client: http.Client{
			Timeout: time.Second * 5,
		},
	}
}
