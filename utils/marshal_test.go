package utils

import (
	demopackage "github.com/go-masonry/mortar/http/server/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalMessageBody(t *testing.T) {
	tests := []struct {
		input  interface{}
		output string
	}{
		{input: []byte(`["marshaled"]`), output: `["marshaled"]`},
		{input: &demopackage.PingRequest{In: "ping"}, output: `{"in": "ping"}`},
		{input: map[string]interface{}{"one": 1}, output: `{"one": 1}`},
	}
	for _, test := range tests {
		bytes, err := MarshalMessageBody(test.input)
		require.NoError(t, err)
		assert.JSONEq(t, test.output, string(bytes))
	}
}
