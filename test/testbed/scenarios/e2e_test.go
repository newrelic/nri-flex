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
	for _, commandTest := range fixtures.CommandTests {
		t.Run(commandTest.Name, func(t *testing.T) {
			tmpConfig, err := tmpFile(commandTest.Config)
			if err != nil {
				t.Error(err)
			}

			defer os.Remove(tmpConfig.Name())

			validator, err := testbed.NewIntegrationValidator(commandTest.ExpectedStdout, "")
			if err != nil {
				t.Error(err)
			}
			tc := testbed.NewTestCase(t, testbed.NewChildFlexRunner("/bin/nri-flex", tmpConfig.Name()), validator)
			tc.RunTest()
		})
	}
}

func TestUrlAPI(t *testing.T) {
	for _, apiTest := range fixtures.URLTests {
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

func TestFileAPI(t *testing.T) {
	for _, urlTest := range fixtures.FileTests {
		t.Run(urlTest.Name, func(t *testing.T) {
			tmpIntegrationFile, err := tmpFile(urlTest.FileContent)
			if err != nil {
				t.Error(err)
			}
			defer os.Remove(tmpIntegrationFile.Name())

			// replace FILE_NAME string of the configuration for the temporal file path
			dynamicConfig := strings.Replace(urlTest.Config, "FILE_PATH", tmpIntegrationFile.Name(), 1)

			tmpConfig, err := tmpFile(dynamicConfig)
			if err != nil {
				t.Error(err)
			}

			defer os.Remove(tmpConfig.Name())

			validator, err := testbed.NewIntegrationValidator(urlTest.ExpectedStdout, "")
			if err != nil {
				t.Error(err)
			}
			tc := testbed.NewTestCase(t, testbed.NewChildFlexRunner("/bin/nri-flex", tmpConfig.Name()), validator)
			tc.RunTest()
		})
	}
}
