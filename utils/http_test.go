package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/go-masonry/mortar/http/client"
	demopackage "github.com/go-masonry/mortar/http/server/proto"
	clientInterface "github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDefaultProtobufHTTPClient(t *testing.T) {
	defaultClient := DefaultProtobufHTTPClient
	if impl, ok := defaultClient.(*protobufHTTPClientImpl); assert.True(t, ok, "wrong implementation") {
		assert.NotNil(t, impl.client)
		assert.NotNil(t, impl.marshaller)
		assert.NotNil(t, impl.errorMapper)
	}
}

func TestDefaultProtobufHTTPClientMethodPostHappy(t *testing.T) {
	var in *demopackage.PingRequest = &demopackage.PingRequest{
		In: "packet",
	}
	method := http.MethodPost
	testDefaultProtobufHTTPClientHappy(t, method, PongTransportWithRequestBodyInterceptor(method), in)
}

func TestDefaultProtobufHTTPClientMethodGetHappy(t *testing.T) {
	method := http.MethodGet
	testDefaultProtobufHTTPClientHappy(t, method, PongTransportWithoutRequestBodyInterceptor(method), nil)
}

func testDefaultProtobufHTTPClientHappy(t *testing.T, method string, interceptor func(req *http.Request, handler clientInterface.HTTPHandler) (*http.Response, error), in proto.Message) {
	client := client.HTTPClientBuilder().AddInterceptors(interceptor).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var out *demopackage.PongResponse
	err := protoClient.Do(context.Background(), method, "http://unreachable", in, &out)
	assert.NoError(t, err)
	assert.Equal(t, "packet", out.GetOut())
}

func PongTransportWithRequestBodyInterceptor(method string) clientInterface.HTTPClientInterceptor {
	return func(req *http.Request, handler clientInterface.HTTPHandler) (*http.Response, error) {
		if method != req.Method {
			return nil, fmt.Errorf("unexpected method: %s", req.Method)
		}

		var bodyMap map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&bodyMap)
		if err != nil {
			return nil, err
		}
		return createHTTPResponse(fmt.Sprintf(`{"out": "%s"}`, bodyMap["in"])), nil
	}
}

func PongTransportWithoutRequestBodyInterceptor(method string) clientInterface.HTTPClientInterceptor {
	return func(req *http.Request, handler clientInterface.HTTPHandler) (*http.Response, error) {
		if method != req.Method {
			return nil, fmt.Errorf("unexpected method: %s", req.Method)
		}

		return createHTTPResponse(`{"out": "packet"}`), nil
	}
}

func createHTTPResponse(httpBody string) *http.Response {
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: 17,
		Body:          ioutil.NopCloser(strings.NewReader(httpBody)),
	}
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

func TestDefaultProtobufHTTPClientsEmptyBodyError(t *testing.T) {
	client := client.HTTPClientBuilder().AddInterceptors(func(_ *http.Request, _ clientInterface.HTTPHandler) (*http.Response, error) {
		return &http.Response{
			Status:        http.StatusText(http.StatusAccepted),
			StatusCode:    http.StatusAccepted,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: 5,
			Body:          ioutil.NopCloser(strings.NewReader("")),
		}, nil
	}).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var in *demopackage.PingRequest = &demopackage.PingRequest{
		In: "packet",
	}
	var empty *emptypb.Empty
	err := protoClient.Do(context.Background(), http.MethodPost, "http://unreachable", in, &empty)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[EOF]")
}

func TestDefaultProtobufHTTPClientIgnoreResponseEvenOnEmptyBody(t *testing.T) {
	client := client.HTTPClientBuilder().AddInterceptors(func(_ *http.Request, _ clientInterface.HTTPHandler) (*http.Response, error) {
		return &http.Response{
			Status:        http.StatusText(http.StatusAccepted),
			StatusCode:    http.StatusAccepted,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: 5,
			Body:          ioutil.NopCloser(strings.NewReader("")),
		}, nil
	}).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var in *demopackage.PingRequest = &demopackage.PingRequest{
		In: "packet",
	}
	err := protoClient.Do(context.Background(), http.MethodPost, "http://unreachable", in, nil)
	assert.NoError(t, err)
}

func TestDefaultProtobufHTTPClientInputAsEmptyBody(t *testing.T) {
	errorMessage := "empty is just fine"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)
		require.EqualValuesf(t, "{}", string(body), "Body: %s", body)
		http.Error(w, errorMessage, http.StatusBadRequest)
	}))
	client := client.HTTPClientBuilder().WithPreconfiguredClient(server.Client()).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var in *emptypb.Empty = &emptypb.Empty{}
	err := protoClient.Do(context.Background(), http.MethodGet, server.URL, in, nil)
	assert.EqualError(t, err, fmt.Sprintf("rpc error: code = InvalidArgument desc = %s\n", errorMessage))
}

func TestDefaultProtobufHTTPClientStatusProtoError(t *testing.T) {
	errorMessage := "this is an error"
	statusErr := status.New(codes.InvalidArgument, errorMessage)
	statusErrJSON, err := protojson.Marshal(statusErr.Proto())
	assert.NoError(t, err)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, string(statusErrJSON), http.StatusBadRequest)
	}))
	client := client.HTTPClientBuilder().WithPreconfiguredClient(server.Client()).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var in *emptypb.Empty = &emptypb.Empty{}
	err = protoClient.Do(context.Background(), http.MethodGet, server.URL, in, nil)
	assert.EqualError(t, err, fmt.Sprintf("rpc error: code = InvalidArgument desc = %s", errorMessage))
}

func TestDefaultProtobufHTTPClientNonStatusProtoJSONError(t *testing.T) {
	errorMessage := `{"key1": "value1", "key2": "value2"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, errorMessage, http.StatusBadRequest)
	}))
	client := client.HTTPClientBuilder().WithPreconfiguredClient(server.Client()).Build()
	protoClient := CreateProtobufHTTPClient(client, nil, nil)
	var in *emptypb.Empty = &emptypb.Empty{}
	err := protoClient.Do(context.Background(), http.MethodGet, server.URL, in, nil)
	assert.EqualError(t, err, fmt.Sprintf("rpc error: code = InvalidArgument desc = %s\n", errorMessage))
}
