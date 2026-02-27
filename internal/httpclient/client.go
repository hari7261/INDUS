package httpclient

import (
	"bytes"
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

// idempotentMethods are safe to retry without side effects.
var idempotentMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodPut:     true,
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	// Read body once so it can be replayed on retries.
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	// Only retry idempotent methods to avoid duplicate side effects.
	maxAttempts := 1
	if idempotentMethods[req.Method] {
		maxAttempts = c.maxRetries
	}

	var (
		resp    *http.Response
		lastErr error
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Check cancellation before each attempt.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Reattach a fresh body reader for every attempt.
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))
		}

		resp, lastErr = c.httpClient.Do(req)

		// Success: 2xx / 3xx / 4xx — caller decides what to do with 4xx.
		if lastErr == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Drain and close the body so the connection can be reused.
		if resp != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			resp = nil
		}

		// No more attempts left.
		if attempt == maxAttempts-1 {
			break
		}

		// Exponential back-off before the next attempt.
		delay := c.baseDelay * time.Duration(1<<uint(attempt))
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed after %d attempt(s): %w", maxAttempts, lastErr)
	}

	// All attempts returned a 5xx status — surface it as an error.
	// resp is nil here because we closed it above, so report via lastErr path never reached;
	// re-issue last attempt's status via a sentinel message.
	return nil, fmt.Errorf("request failed after %d attempt(s): server returned 5xx", maxAttempts)
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
