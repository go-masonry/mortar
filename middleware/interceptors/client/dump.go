package client

import (
	"net/http"
	"net/http/httputil"

	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

type dumpHTTPDeps struct {
	fx.In

	Logger log.Logger
}

// DumpRESTClientInterceptor usefull when you want to log what is actually sent to the external HTTP server
// and was returned.
func DumpRESTClientInterceptor(deps dumpHTTPDeps) client.HTTPClientInterceptor {
	return func(req *http.Request, handler client.HTTPHandler) (*http.Response, error) {
		reqBody, err := httputil.DumpRequestOut(req, true)
		deps.Logger.WithError(err).Debug(req.Context(), "Request:\n%s\n", reqBody)
		res, err := handler(req)
		if err == nil {
			resBody, dumpErr := httputil.DumpResponse(res, true)
			deps.Logger.WithError(dumpErr).Debug(req.Context(), "Response:\n%s\n", resBody)
		}
		return res, err
	}
}
