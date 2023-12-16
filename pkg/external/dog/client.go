package dog

import (
	"context"
	"fmt"
	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
	"io"
	"net/http"
	"time"
)

type ClientConfig struct {
	Enabled                bool
	BaseUrl                string
	CircuitBreakerSettings gobreaker.Settings
}

type Client struct {
	ClientConfig
	CircuitBreaker *gobreaker.CircuitBreaker
	HttpClient     *http.Client
}

type NoOpClient struct {
}

func NewClient(config ClientConfig) user.Dog {
	if !config.Enabled {
		return NoOpClient{}
	}
	return &Client{
		ClientConfig:   config,
		CircuitBreaker: gobreaker.NewCircuitBreaker(config.CircuitBreakerSettings),
		HttpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) GetRandomDogUrl(ctx context.Context) (string, error) {
	l := zerolog.Ctx(ctx)

	path := fmt.Sprintf(c.BaseUrl + "/woof.json")
	l.Debug().Msgf("GetRandomDogUrl: path: %s", path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", fmt.Errorf("GetRandomDogUrl: could not create request: %w", err)
	}

	r, err := c.CircuitBreaker.Execute(func() (interface{}, error) {
		return c.HttpClient.Do(req)
	})
	if err != nil {
		return "", fmt.Errorf("GetRandomDogUrl: could not execute GET: %w", err)
	}

	resp := r.(*http.Response)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("GetRandomDogUrl: could not read response")
	}

	url := struct {
		Url string `json:"url"`
	}{}
	if err := json.Unmarshal(bytes, &url); err != nil {
		return "", fmt.Errorf("GetRandomDogUrl: body: %s could not decode response: %w", string(bytes), err)
	}

	return url.Url, nil
}

func (c NoOpClient) GetRandomDogUrl(context.Context) (string, error) {
	return "", nil
}
