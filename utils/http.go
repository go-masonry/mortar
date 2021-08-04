package utils

import (
	"bytes"
	"context"
	"net/http"
	"reflect"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ErrorMapper is a map function that maps HTTP Status Code into its gRPC counter part.
type ErrorMapper func(statusCode int) *status.Status

// ProtobufHTTPClient is a helper util in situations where you want to call a REST API, but you have all the definitions as Protobuf.
type ProtobufHTTPClient interface {
	// Example:
	//	var response *pbpkg.ResponseMessage
	//	var request = &pbpkg.RequestMessage{Name: "test"}
	//	err := Do(ctx, http.MethodPost, "http://host/path", request, &response)
	//
	// Error returned will be of type `(*status.Status).Err()` meaning it will be a gRPC type error.
	//
	// Note:
	// you can pass `nil` as the last parameter if you don't want to unmarshal HTTP response body.
	// or you know it's going to be empty == EOF
	//	err := Do(ctx, http.MethodPost, "http://host/path", request, nil)
	Do(ctx context.Context, method, url string, in proto.Message, out interface{}) error
}

// DefaultProtobufHTTPClient uses default http Client, error mapper and marshaller
var DefaultProtobufHTTPClient = CreateProtobufHTTPClient(nil, nil, nil)

// CreateProtobufHTTPClient Creates a custom Protobuf aware HTTP client
func CreateProtobufHTTPClient(client *http.Client, errorMapper ErrorMapper, marshaller runtime.Marshaler) ProtobufHTTPClient {
	if client == nil {
		client = http.DefaultClient
	}
	if errorMapper == nil {
		errorMapper = defaultErrorMapper
	}
	if marshaller == nil {
		marshaller = &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}
	}
	return &protobufHTTPClientImpl{
		client:      client,
		errorMapper: errorMapper,
		marshaller:  marshaller,
	}
}

type protobufHTTPClientImpl struct {
	client      *http.Client
	errorMapper ErrorMapper
	marshaller  runtime.Marshaler
}

// ProtoToHTTPRequest is a helper to convert proto Message into an HTTP Request
func (impl *protobufHTTPClientImpl) Do(ctx context.Context, method, url string, in proto.Message, out interface{}) error {
	var reqAsBytes []byte
	var response *http.Response
	var request *http.Request
	var err error
	reqAsBytes, err = impl.marshaller.Marshal(in)
	if err != nil {
		return status.Errorf(codes.Unknown, "error while marshaling request, %s", err)
	}
	buffer := bytes.NewBuffer(reqAsBytes)
	request, err = http.NewRequestWithContext(ctx, method, url, buffer)
	if err != nil {
		return status.Errorf(codes.Unknown, "error while creating an http request, %s", err)
	}
	response, err = impl.client.Do(request)
	if err != nil {
		return status.Errorf(codes.Internal, "error executing http call, %s", err)
	}
	defer response.Body.Close()
	if grpcStatus := impl.errorMapper(response.StatusCode); grpcStatus != nil && grpcStatus.Code() != codes.OK {
		var responseBodyBuffer bytes.Buffer
		if _, err = responseBodyBuffer.ReadFrom(response.Body); err == nil {
			var statusError *spb.Status
			responseBodyBytes := responseBodyBuffer.Bytes()
			if decodeError := impl.marshaller.NewDecoder(bytes.NewReader(responseBodyBytes)).Decode(&statusError); decodeError == nil {
				if statusError.GetCode() != int32(codes.OK) {
					return status.ErrorProto(statusError)
				}
			}
			return status.Error(grpcStatus.Code(), string(responseBodyBytes))
		}
		return grpcStatus.Err()
	}
	// no need to unmarshal the body if it's a "nil interface"
	// https://golang.org/doc/faq#nil_error
	if reflect.TypeOf(out) != nil {
		if decodeError := impl.marshaller.NewDecoder(response.Body).Decode(out); decodeError != nil {
			return status.Errorf(codes.Unknown, "error unmarshalling response, [%s]", decodeError)
		}
	}
	return nil
}

func defaultErrorMapper(httpStatus int) *status.Status {
	switch httpStatus {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return status.New(codes.OK, http.StatusText(httpStatus))
	case http.StatusBadRequest:
		return status.New(codes.InvalidArgument, http.StatusText(httpStatus))
	case http.StatusMethodNotAllowed:
		return status.New(codes.Unimplemented, http.StatusText(httpStatus))
	case http.StatusNotFound:
		return status.New(codes.NotFound, http.StatusText(httpStatus))
	case http.StatusConflict:
		return status.New(codes.AlreadyExists, http.StatusText(httpStatus))
	case http.StatusUnauthorized:
		return status.New(codes.Unauthenticated, http.StatusText(httpStatus))
	case http.StatusTooManyRequests:
		return status.New(codes.ResourceExhausted, http.StatusText(httpStatus))
	case http.StatusNotImplemented:
		return status.New(codes.Unimplemented, http.StatusText(httpStatus))
	case http.StatusInternalServerError:
		return status.New(codes.Internal, http.StatusText(httpStatus))
	case http.StatusServiceUnavailable:
		return status.New(codes.Unavailable, http.StatusText(httpStatus))
	default:
		return status.New(codes.Unknown, http.StatusText(httpStatus))
	}
}
