package outputs

import (
	"compress/zlib"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func Test_postRequest_no_such_host(t *testing.T) {
	someData := []byte("some data")
	someKey := "key"
	load.Refresh()
	load.Entity = &integration.Entity{
		Metrics: []*metric.Set{},
	}
	err := postRequest("http://bad..url..z/", someKey, someData)
	require.EqualError(t, err, "http: failed to send: Post \"http://bad..url..z/\": dial tcp: lookup bad..url..z: no such host")
}

func Test_postRequest_create_request_fail(t *testing.T) {
	someData := []byte("some data")
	someKey := "key"
	load.Refresh()
	load.Entity = &integration.Entity{
		Metrics: []*metric.Set{},
	}
	err := postRequest("%zzzzz", someKey, someData)
	require.EqualError(t, err, "http: unable to create http.Request, parse \"%zzzzz\": invalid URL escape \"%zz\"")
}

func Test_postRequest_http_post(t *testing.T) {
	someData := []byte("some data")
	someKey := "key"
	load.Refresh()
	load.Entity = &integration.Entity{
		Metrics: []*metric.Set{},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var output string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		_, _ = fmt.Fprintln(w, "Hello, client")
		reader, err := zlib.NewReader(r.Body)
		require.NoError(t, err)
		b, err := ioutil.ReadAll(reader)
		require.NoError(t, err)
		output = string(b)
	}))

	defer ts.Close()
	err := postRequest(ts.URL, someKey, someData)
	wg.Wait()
	require.NoError(t, err)
	assert.Equal(t, "some data", output)
}

func Test_postRequest_http_post_300(t *testing.T) {
	someData := []byte("some data")
	someKey := "key"
	load.Refresh()
	load.Entity = &integration.Entity{
		Metrics: []*metric.Set{},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		w.WriteHeader(300)
		_, _ = ioutil.ReadAll(r.Body)
	}))

	defer ts.Close()
	err := postRequest(ts.URL, someKey, someData)
	wg.Wait()
	require.EqualError(t, err, "http: post failed, status code: 300")
}
