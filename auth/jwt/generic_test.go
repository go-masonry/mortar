package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const fakeToken = "fakeAlg.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.fakeSignature"
const fakeTokenBody = `{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}`

type body struct {
	Sub  string  `json:"sub"`
	Name string  `json:"name"`
	Iat  float64 `json:"iat"`
}

func TestDefaults(t *testing.T) {
	extractor := Builder().Build()
	token, err := extractor.FromString(fakeToken)
	assert.NoError(t, err)
	assert.JSONEq(t, fakeTokenBody, string(token.Payload()))
	_, err = extractor.FromContext(context.Background())
	assert.Error(t, err)
}

func TestCustom(t *testing.T) {
	extractor := Builder().SetBase64Decoder(base64.RawURLEncoding).SetDecoder(json.Unmarshal).Build()
	token, err := extractor.FromString(fakeToken)
	assert.NoError(t, err)
	assert.JSONEq(t, fakeTokenBody, string(token.Payload()))
	var b body
	err = token.Decode(&b)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", b.Name)
}

func TestContextExtractor(t *testing.T) {
	ctxExtractor := func(ctx context.Context) (string, error) {
		return fakeToken, nil
	}
	extractor := Builder().SetContextExtractor(ctxExtractor).Build()
	token, err := extractor.FromContext(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, fakeToken, token.Raw())
}

func TestTokenToMap(t *testing.T) {
	extractor := Builder().Build()
	token, err := extractor.FromString(fakeToken)
	assert.NoError(t, err)
	expectedMap := map[string]interface{}{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  float64(1516239022),
	}
	actual, err := token.Map()
	assert.NoError(t, err)
	assert.Equal(t, expectedMap, actual)
}

func TestTokenToStruct(t *testing.T) {
	extractor := Builder().Build()
	token, err := extractor.FromString(fakeToken)
	assert.NoError(t, err)
	var expected = body{
		Sub:  "1234567890",
		Name: "John Doe",
		Iat:  float64(1516239022),
	}
	var actual body
	err = token.Decode(&actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBadToken(t *testing.T) {
	extractor := Builder().Build()
	_, err := extractor.FromString("fake string num 1")
	assert.Error(t, err)
	_, err = extractor.FromString("part1.part2.part3")
	assert.Error(t, err)
}
