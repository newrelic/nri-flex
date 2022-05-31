package scenarios

import (
	"github.com/newrelic/nri-flex/test/testbed"
	"github.com/newrelic/nri-flex/test/testbed/scenarios/fixtures"
	"io/ioutil"
	"os"
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

func TestDiskLinux(t *testing.T) {
	for _, diskTest := range fixtures.DiskTests {
		t.Run(diskTest.Name, func(t *testing.T) {
			tmpConfig, err := tmpFile(diskTest.Config)
			if err != nil {
				t.Error(err)
			}

			defer os.Remove(tmpConfig.Name())

			validator, err := testbed.NewMetricValidator(diskTest.ExpectedStdout, "nothing")
			if err != nil {
				t.Error(err)
			}
			tc := testbed.NewTestCase(t, testbed.NewChildFlexRunner("/bin/nri-flex", tmpConfig.Name()), validator)
			tc.RunTest()
		})
	}
}
