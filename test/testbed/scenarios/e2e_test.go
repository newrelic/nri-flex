package scenarios

import (
	"github.com/newrelic/nri-flex/test/testbed"
	"github.com/newrelic/nri-flex/test/testbed/scenarios/fixtures"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func tmpFile(data string) (file *os.File, err error) {
	file, err = ioutil.TempFile("", "")
	if err != nil {
		return
	}
	_, err = file.Write([]byte(data))
	file.Close()
	return
}

func TestCommandAPI(t *testing.T) {
	for _, diskTest := range fixtures.CommandTests {
		t.Run(diskTest.Name, func(t *testing.T) {
			tmpConfig, err := tmpFile(diskTest.Config)
			if err != nil {
				t.Error(err)
			}

			defer os.Remove(tmpConfig.Name())

			validator, err := testbed.NewIntegrationValidator(diskTest.ExpectedStdout, "nothing")
			if err != nil {
				t.Error(err)
			}
			tc := testbed.NewTestCase(t, testbed.NewChildFlexRunner("/bin/nri-flex", tmpConfig.Name()), validator)
			tc.RunTest()
		})
	}
}

func TestAPI(t *testing.T) {
	for _, apiTest := range fixtures.UrlTests {
		t.Run(apiTest.Name, func(t *testing.T) {

			tmpConfig, err := tmpFile(apiTest.Config)
			if err != nil {
				t.Error(err)
			}

			defer os.Remove(tmpConfig.Name())

			srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "/"+apiTest.Endpoint) {
					_, err := w.Write([]byte(apiTest.Payload))
					assert.NoError(t, err)
				} else {
					_, err := w.Write([]byte{})
					assert.NoError(t, err)
				}
			}))
			l, _ := net.Listen("tcp", "127.0.0.1:"+apiTest.Port)
			srv.Listener = l
			srv.Start()
			defer srv.Close()
			validator, err := testbed.NewIntegrationValidator(apiTest.ExpectedStdout, "nothing")
			if err != nil {
				t.Error(err)
			}
			tc := testbed.NewTestCase(t, testbed.NewChildFlexRunner("/bin/nri-flex", tmpConfig.Name()), validator)
			tc.RunTest()
		})
	}
}
