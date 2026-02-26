package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	maxRetries int
	baseDelay  time.Duration
}

type Options struct {
	Timeout    time.Duration
	MaxRetries int
}

func New(opts Options) *Client {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	
	return &Client{
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		maxRetries: opts.MaxRetries,
		baseDelay:  1 * time.Second,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	
	var resp *http.Response
	var err error
	
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		// Check cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		resp, err = c.httpClient.Do(req)
		
		// Success
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		
		// Close response body if present
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		
		// Don't retry on last attempt
		if attempt == c.maxRetries-1 {
			break
		}
		
		// Exponential backoff
		delay := c.baseDelay * time.Duration(1<<uint(attempt))
		
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries, err)
	}
	
	return resp, nil
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}
