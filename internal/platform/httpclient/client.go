package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client is a centralized HTTP client with built-in observability, logging, and middleware support
type Client struct {
	httpClient     *http.Client
	logger         *zap.Logger
	interceptors   []Interceptor
	headerProvider HeaderProvider
	serviceName    string
}

// Config holds configuration for the HTTP client
type Config struct {
	Timeout        time.Duration
	MaxIdleConns   int
	ServiceName    string
	HeaderProvider HeaderProvider // Optional: provides headers for all requests
}

// NewClient creates a new instrumented HTTP client
func NewClient(cfg Config, logger *zap.Logger, interceptors ...Interceptor) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 100
	}

	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConns,
		IdleConnTimeout:     90 * time.Second,
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		logger:         logger,
		interceptors:   interceptors,
		headerProvider: cfg.HeaderProvider,
		serviceName:    cfg.ServiceName,
	}
}

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}

// Response represents an HTTP response with metadata
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

// Do executes an HTTP request with full observability
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	startTime := time.Now()

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, req.Body)
	if err != nil {
		c.logger.Error("failed to create request",
			zap.String("service", c.serviceName),
			zap.String("method", req.Method),
			zap.String("url", req.URL),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers from provider (if configured)
	if c.headerProvider != nil {
		providedHeaders, err := c.headerProvider.GetHeaders(ctx)
		if err != nil {
			c.logger.Error("failed to get headers from provider",
				zap.String("service", c.serviceName),
				zap.Error(err),
			)
			return nil, fmt.Errorf("header provider failed: %w", err)
		}
		for key, value := range providedHeaders {
			httpReq.Header.Set(key, value)
		}
	}

	// Set request-specific headers (can override provider headers)
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute interceptors (before)
	for _, interceptor := range c.interceptors {
		if err := interceptor.Before(ctx, httpReq); err != nil {
			return nil, fmt.Errorf("interceptor before failed: %w", err)
		}
	}

	// Log request
	c.logger.Info("http request initiated",
		zap.String("service", c.serviceName),
		zap.String("method", req.Method),
		zap.String("url", req.URL),
	)

	// Execute request
	httpResp, err := c.httpClient.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error("http request failed",
			zap.String("service", c.serviceName),
			zap.String("method", req.Method),
			zap.String("url", req.URL),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		c.logger.Error("failed to read response body",
			zap.String("service", c.serviceName),
			zap.String("method", req.Method),
			zap.String("url", req.URL),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Body:       body,
		Duration:   duration,
	}

	// Execute interceptors (after)
	for _, interceptor := range c.interceptors {
		if err := interceptor.After(ctx, httpResp, response); err != nil {
			c.logger.Warn("interceptor after failed",
				zap.String("service", c.serviceName),
				zap.Error(err),
			)
		}
	}

	// Log response with full observability
	logLevel := zap.InfoLevel
	if httpResp.StatusCode >= 400 {
		logLevel = zap.ErrorLevel
	}

	c.logger.Log(logLevel, "http request completed",
		zap.String("service", c.serviceName),
		zap.String("method", req.Method),
		zap.String("url", req.URL),
		zap.Int("status_code", httpResp.StatusCode),
		zap.Duration("duration", duration),
		zap.Int("response_size", len(body)),
	)

	return response, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
	})
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:  http.MethodPost,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, url string, body io.Reader, headers map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:  http.MethodPut,
		URL:     url,
		Body:    body,
		Headers: headers,
	})
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
	})
}

// GetJSON performs a GET request and unmarshals JSON response
func (c *Client) GetJSON(ctx context.Context, url string, result interface{}) error {
	resp, err := c.Get(ctx, url, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: request failed", resp.StatusCode)
	}

	if result != nil {
		if err := json.Unmarshal(resp.Body, result); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	return nil
}

// PostJSON performs a POST request with JSON body and unmarshals JSON response
func (c *Client) PostJSON(ctx context.Context, url string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	resp, err := c.Post(ctx, url, bodyReader, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	})
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: request failed", resp.StatusCode)
	}

	if result != nil {
		if err := json.Unmarshal(resp.Body, result); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	return nil
}
