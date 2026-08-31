package geminiclient

import (
	"context"

	httpx "github.com/AppeiYA/consultation-platform/internal/shared/adapters/http/client"
	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	model string
}

func NewClient(ctx context.Context, httpClient *httpx.Client, config Config) (*Client, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.APIKey,
		Backend: genai.BackendGeminiAPI,
		HTTPClient: httpClient.Client,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		model:  config.Model,
	}, nil
}

func (c *Client) Model() string {
	return c.model
}

func (c *Client) Client() *genai.Client {
	return c.client
}

func (c *Client) GenerateContent(
	ctx context.Context, 
	contents []*genai.Content, 
	config *genai.GenerateContentConfig,
) (*genai.GenerateContentResponse, error) {
	return c.client.Models.GenerateContent(
		ctx,
		c.model,
		contents,
		config,
	)
}