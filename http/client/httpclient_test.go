package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/stretchr/testify/require"
)

func TestCustomClient(t *testing.T) {
	expected := &http.Client{}
	actual := HTTPClientBuilder().WithPreconfiguredClient(expected).Build()
	require.Equal(t, expected, actual, "other client")
}

func TestDefault(t *testing.T) {
	client := HTTPClientBuilder().Build()
	require.NotNil(t, client, "an empty client")
}

func TestWithInterceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		name := req.Header.Get("name")
		_, err := writer.Write([]byte("Hello " + name))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := HTTPClientBuilder().WithPreconfiguredClient(server.Client()).AddInterceptors(testInterceptor).Build()
	response, err := client.Get(server.URL)
	require.NoError(t, err)
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Contains(t, string(bodyBytes), "Hello Robert")
	family := response.Header.Get("family")
	require.Equal(t, "Pike", family)
}

func testInterceptor(req *http.Request, handler client.HTTPHandler) (resp *http.Response, err error) {
	req.Header.Set("name", "Robert")
	if resp, err = handler(req); err == nil {
		resp.Header.Set("family", "Pike")
	}
	return
}
