package httpclient

import (
	"context"
	"net/http"
)

// Interceptor defines the interface for request/response interceptors
type Interceptor interface {
	// Before is called before the request is sent
	Before(ctx context.Context, req *http.Request) error

	// After is called after the response is received
	After(ctx context.Context, resp *http.Response, response *Response) error
}

// InterceptorFunc is a function type that implements the Interceptor interface
type InterceptorFunc struct {
	BeforeFunc func(ctx context.Context, req *http.Request) error
	AfterFunc  func(ctx context.Context, resp *http.Response, response *Response) error
}

func (i InterceptorFunc) Before(ctx context.Context, req *http.Request) error {
	if i.BeforeFunc != nil {
		return i.BeforeFunc(ctx, req)
	}
	return nil
}

func (i InterceptorFunc) After(ctx context.Context, resp *http.Response, response *Response) error {
	if i.AfterFunc != nil {
		return i.AfterFunc(ctx, resp, response)
	}
	return nil
}
