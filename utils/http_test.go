package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/go-masonry/mortar/http/client"
	demopackage "github.com/go-masonry/mortar/http/server/proto"
	clientInterface "github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDefaultProtobufHTTPClient(t *testing.T) {
	defaultClient := DefaultProtobufHTTPClient
	if impl, ok := defaultClient.(*protobufHTTPClientImpl); assert.True(t, ok, "wrong implementation") {
		assert.NotNil(t, impl.client)
		assert.NotNil(t, impl.marshaller)
		assert.NotNil(t, impl.errorMapper)
	}
}

func TestDefaultProtobufHTTPClientHappy(t *testing.T) {
	client := client.HTTPClientBuilder().AddInterceptors(PongTransportInterceptor).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var in *demopackage.PingRequest = &demopackage.PingRequest{
		In: "packet",
	}
	var out *demopackage.PongResponse
	err := protoClient.Do(context.Background(), http.MethodPost, "http://unreachable", in, &out)
	assert.NoError(t, err)
	assert.Equal(t, "packet", out.GetOut())
}

func PongTransportInterceptor(req *http.Request, handler clientInterface.HTTPHandler) (*http.Response, error) {
	var bodyMap map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&bodyMap)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 17,
		Body:          ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"out": "%s"}`, bodyMap["in"]))),
	}, nil
}

func TestDefaultProtobufHTTPClientErrorMapping(t *testing.T) {
	httpErrors := map[int]codes.Code{
		http.StatusBadRequest:          codes.InvalidArgument,
		http.StatusMethodNotAllowed:    codes.Unimplemented,
		http.StatusNotFound:            codes.NotFound,
		http.StatusConflict:            codes.AlreadyExists,
		http.StatusUnauthorized:        codes.Unauthenticated,
		http.StatusTooManyRequests:     codes.ResourceExhausted,
		http.StatusNotImplemented:      codes.Unimplemented,
		http.StatusInternalServerError: codes.Internal,
		http.StatusServiceUnavailable:  codes.Unavailable,
		http.StatusTeapot:              codes.Unknown,
	}
	for httpCode, grpcCode := range httpErrors {
		client := client.HTTPClientBuilder().AddInterceptors(func(_ *http.Request, _ clientInterface.HTTPHandler) (*http.Response, error) {
			return &http.Response{
				Status:        http.StatusText(httpCode),
				StatusCode:    httpCode,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: 5,
				Body:          ioutil.NopCloser(strings.NewReader("error")),
			}, nil
		}).Build()
		protoClient := CreateProtobufHTTPClient(client, nil, nil)
		var in *demopackage.PingRequest = &demopackage.PingRequest{
			In: "packet",
		}
		var out *demopackage.PongResponse
		err := protoClient.Do(context.Background(), http.MethodPost, "http://unreachable", in, &out)
		assert.Error(t, err)
		assert.Nil(t, out)
		require.Implements(t, (*(interface{ GRPCStatus() *status.Status }))(nil), err)
		sts := err.(interface{ GRPCStatus() *status.Status })
		assert.Equal(t, grpcCode.String(), sts.GRPCStatus().Code().String(), sts.GRPCStatus().Message())
	}
}
