package jwt

import (
	"context"
	"encoding/base64"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// JSONDecoder is an alias to json.Unmarshal
type JSONDecoder func(data []byte, v interface{}) error

// ExtractorBuilder lets you define custom options
type ExtractorBuilder interface {
	// What decoder to use when unmarshalling
	SetDecoder(dec JSONDecoder) ExtractorBuilder
	// ContextExtractor lets you set a custom extractor from context.Context
	SetContextExtractor(extractor ContextExtractor) ExtractorBuilder
	// SetBase64Decoder lets you customize base64.Encoding, standard or URL or other
	SetBase64Decoder(dec *base64.Encoding) ExtractorBuilder
	Build() TokenExtractor
}

// ContextExtractor is a helper function used to extract values from context
// Sometimes it's convenient to have Authorization header within the context, i.e GRPC-Gateway
// Implementation of this function should treat that use case
type ContextExtractor func(ctx context.Context) (string, error)

// TokenExtractor is a public interface to help with token extraction or preparing
type TokenExtractor interface {
	// FromContext should try to extract a token from the Context using `ContextExtractor`
	FromContext(ctx context.Context) (Token, error)
	// FromString accepts a JWT token in form of a string
	//	xxxxx.yyyyy.zzzzz
	FromString(str string) (Token, error)
}

// Token interface shouldn't be used as a standalone since it doesn't have anything to work on
// It's implementation must be created by the `Extractor` interface from above
type Token interface {
	// Raw returns JWT token as is: Base64("<algo>.<payload>.<signature>")
	Raw() string
	// Payload returns the payload (middle part) of JWT after base64 decode: UnBase64("<payload>")
	Payload() []byte
	// Map extracts token values (it's a JSON) to a map
	Map() (map[string]interface{}, error)
	// Decode extracts token values to a specified struct, struct should be a pointer
	// 	{
	// 		Subject string `json:"sub"`
	// 		Issuer string `json:"iss"`
	//		...
	// 	}
	Decode(target interface{}) error
}
