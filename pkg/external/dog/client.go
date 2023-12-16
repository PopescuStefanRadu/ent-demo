package dog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/PopescuStefanRadu/ent-demo/pkg/user"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/sony/gobreaker"
)

var ErrCouldNotReadResponse = errors.New("GetRandomDogURL: could not read response")

type ClientConfig struct {
	Enabled                bool
	BaseURL                string
	CircuitBreakerSettings gobreaker.Settings
}

type Client struct {
	ClientConfig
	CircuitBreaker *gobreaker.CircuitBreaker
	HTTPClient     *http.Client
}

type NoOpClient struct{}

//nolint:ireturn,nolintlint
func NewClient(config ClientConfig) user.Dog {
	if !config.Enabled {
		return NoOpClient{}
	}

	return &Client{
		ClientConfig:   config,
		CircuitBreaker: gobreaker.NewCircuitBreaker(config.CircuitBreakerSettings),
		HTTPClient:     &http.Client{},
	}
}

func (c *Client) GetRandomDogURL(ctx context.Context) (string, error) {
	l := zerolog.Ctx(ctx)

	path := fmt.Sprintf(c.BaseURL + "/woof.json")
	l.Debug().Msgf("GetRandomDogURL: path: %s", path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", fmt.Errorf("GetRandomDogURL: could not create request: %w", err)
	}

	r, err := c.CircuitBreaker.Execute(func() (interface{}, error) {
		return c.HTTPClient.Do(req) //nolint:bodyclose
	})
	if err != nil {
		return "", fmt.Errorf("GetRandomDogURL: could not execute GET: %w", err)
	}

	resp := r.(*http.Response) //nolint:forcetypeassert
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrCouldNotReadResponse
	}

	url := struct {
		URL string `json:"url"`
	}{}
	if err := json.Unmarshal(bytes, &url); err != nil {
		return "", fmt.Errorf("GetRandomDogURL: body: %s could not decode response: %w", string(bytes), err)
	}

	return url.URL, nil
}

func (c NoOpClient) GetRandomDogURL(context.Context) (string, error) {
	return "", nil
}
