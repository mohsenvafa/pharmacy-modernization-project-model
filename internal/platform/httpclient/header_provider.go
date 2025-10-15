package httpclient

import "context"

// HeaderProvider is an interface for providing dynamic headers
type HeaderProvider interface {
	GetHeaders(ctx context.Context) (map[string]string, error)
}

// StaticHeaderProvider provides static headers
type StaticHeaderProvider struct {
	headers map[string]string
}

// NewStaticHeaderProvider creates a provider with static headers
func NewStaticHeaderProvider(headers map[string]string) *StaticHeaderProvider {
	return &StaticHeaderProvider{
		headers: headers,
	}
}

func (p *StaticHeaderProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	return p.headers, nil
}

// HeaderProviderFunc is a function type that implements HeaderProvider
type HeaderProviderFunc func(ctx context.Context) (map[string]string, error)

func (f HeaderProviderFunc) GetHeaders(ctx context.Context) (map[string]string, error) {
	return f(ctx)
}
