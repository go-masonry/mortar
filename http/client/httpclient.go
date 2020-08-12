package client

import (
	"container/list"
	"net/http"

	"github.com/go-masonry/mortar/interfaces/http/client"
)

type restBuilderConfig struct {
	predefinedClient *http.Client
	interceptors     []client.HTTPClientInterceptor
}

type builderImpl struct {
	ll *list.List
}

// HTTPClientBuilder creates a fresh http.Client builder
func HTTPClientBuilder() client.HTTPClientBuilder {
	return &builderImpl{
		ll: list.New(),
	}
}

func (impl *builderImpl) AddInterceptors(interceptors ...client.HTTPClientInterceptor) client.HTTPClientBuilder {
	impl.ll.PushBack(func(cfg *restBuilderConfig) {
		if len(interceptors) > 0 {
			cfg.interceptors = append(cfg.interceptors, interceptors...)
		}
	})
	return impl
}

func (impl *builderImpl) WithPreconfiguredClient(client *http.Client) client.HTTPClientBuilder {
	impl.ll.PushBack(func(cfg *restBuilderConfig) {
		cfg.predefinedClient = client
	})
	return impl
}

func (impl *builderImpl) Build() *http.Client {
	var client = &http.Client{}
	if impl != nil {
		cfg := new(restBuilderConfig)
		for e := impl.ll.Front(); e != nil; e = e.Next() {
			f := e.Value.(func(cfg *restBuilderConfig))
			f(cfg)
		}

		if cfg.predefinedClient != nil {
			client = cfg.predefinedClient
		}
		if client.Transport == nil {
			client.Transport = http.DefaultTransport
		}
		client.Transport = prepareCustomRoundTripper(client.Transport, cfg.interceptors...)
	}

	return client
}

type customRoundTripper struct {
	inner             http.RoundTripper
	unitedInterceptor client.HTTPClientInterceptor
}

func prepareCustomRoundTripper(actual http.RoundTripper, interceptors ...client.HTTPClientInterceptor) http.RoundTripper {
	return &customRoundTripper{
		inner:             actual,
		unitedInterceptor: uniteInterceptors(interceptors),
	}
}

func (crt *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return crt.unitedInterceptor(req, crt.inner.RoundTrip)
}

func uniteInterceptors(interceptors []client.HTTPClientInterceptor) client.HTTPClientInterceptor {
	if len(interceptors) == 0 {
		return func(req *http.Request, handler client.HTTPpHandler) (*http.Response, error) {
			// That's why we needed an alias to http.RoundTripper.RoundTrip
			return handler(req)
		}
	}

	return func(req *http.Request, handler client.HTTPpHandler) (*http.Response, error) {
		tailHandler := func(innerReq *http.Request) (*http.Response, error) {
			unitedInterceptor := uniteInterceptors(interceptors[1:])
			return unitedInterceptor(req, handler)
		}
		headInterceptor := interceptors[0]
		return headInterceptor(req, tailHandler)
	}
}

var _ client.HTTPClientBuilder = (*builderImpl)(nil)
