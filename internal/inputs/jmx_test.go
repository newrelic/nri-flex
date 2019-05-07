package inputs

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestSetJMXCommand(t *testing.T) {
	load.Refresh()
	configs := []load.Config{
		{
			Name: "jmxFlex",
			Global: load.Global{
				Jmx: load.JMX{
					Host:           "127.0.0.1",
					Port:           "9001",
					User:           "batman",
					Pass:           "robin",
					KeyStore:       "abc",
					KeyStorePass:   "def",
					TrustStore:     "abc",
					TrustStorePass: "def",
				},
			},
			APIs: []load.API{
				{
					Name: "tomcat",
					Commands: []load.Command{
						{
							Output: "jmx",
							Run:    "Catalina:type=ThreadPool,name=*",
						},
					},
				},
			},
		},
		{
			Name: "jmxFlex",
			APIs: []load.API{
				{
					Name: "tomcat",
					Jmx: load.JMX{
						Host:           "127.0.0.1",
						Port:           "9001",
						User:           "batman",
						Pass:           "robin",
						KeyStore:       "abc",
						KeyStorePass:   "def",
						TrustStore:     "abc",
						TrustStorePass: "def",
					},
					Commands: []load.Command{
						{
							Output: "jmx",
							Run:    "Catalina:type=ThreadPool,name=*",
						},
					},
				},
			},
		},
		{
			Name: "jmxFlex",
			APIs: []load.API{
				{
					Name: "tomcat",
					Commands: []load.Command{
						{
							Output: "jmx",
							Run:    "Catalina:type=ThreadPool,name=*",
							Jmx: load.JMX{
								Host:           "127.0.0.1",
								Port:           "9001",
								User:           "batman",
								Pass:           "robin",
								KeyStore:       "abc",
								KeyStorePass:   "def",
								TrustStore:     "abc",
								TrustStorePass: "def",
							},
						},
					},
				},
			},
		},
		{
			Name: "jmxFlex",
			APIs: []load.API{
				{
					Name: "tomcat",
					Commands: []load.Command{
						{
							Output: "jmx",
							Run:    "Catalina:type=ThreadPool,name=*",
							Jmx: load.JMX{
								Port:           "9001",
								User:           "batman",
								Pass:           "robin",
								KeyStore:       "abc",
								KeyStorePass:   "def",
								TrustStore:     "abc",
								TrustStorePass: "def",
							},
						},
					},
				},
			},
		},
	}

	for _, config := range configs {
		runCommand := config.APIs[0].Commands[0].Run
		SetJMXCommand(&runCommand, config.APIs[0].Commands[0], config.APIs[0], &config)
		expectedString := `echo "Catalina:type=ThreadPool,name=*" | java -jar ./nrjmx/nrjmx.jar` +
			` -hostname 127.0.0.1 -port 9001 -username batman -password robin ` +
			`-keyStore abc -keyStorePassword def -trustStore abc -trustStorePassword def`
		if runCommand != expectedString {
			t.Errorf("want: %v, got: %v", expectedString, runCommand)
		}
	}

}

func TestParseJMX(t *testing.T) {
	b, _ := ioutil.ReadFile("../../test/payloads/tomcatJMX.out")
	_, dataInterface := detectCommandOutput(string(b), "jmx")
	prevSample := map[string]interface{}{
		"hi": "hello",
	}
	config := load.Config{
		Name: "jmxFlex",
		APIs: []load.API{
			{
				Name: "tomcat",
				Commands: []load.Command{
					{
						Output: "jmx",
						Run:    "Catalina:type=ThreadPool,name=*",
						Jmx: load.JMX{
							Port:           "9001",
							User:           "batman",
							Pass:           "robin",
							KeyStore:       "abc",
							KeyStorePass:   "def",
							TrustStore:     "abc",
							TrustStorePass: "def",
						},
					},
				},
			},
		},
	}
	expectedDatastore := []interface{}{
		map[string]interface{}{
			"Catalina:type":          "ThreadPool",
			"acceptorThreadCount":    1,
			"acceptorThreadPriority": 5,
			"algorithm":              "SunX509",
			"backlog":                100,
			"bean":                   "type=ThreadPool,name=\"http-apr-8080\"",
			"bindOnInit":             true,
			"clientAuth":             "false",
			"connectionCount":        1,
			"currentThreadCount":     10,
			"currentThreadsBusy":     0,
			"daemon":                 true,
			"deferAccept":            true,
			"domain":                 "Catalina",
			"executorTerminationTimeoutMillis": 5000,
			"hi":                         "hello",
			"keepAliveCount":             0,
			"keepAliveTimeout":           20000,
			"keystoreFile":               "/root/.keystore",
			"keystoreType":               "JKS",
			"localPort":                  8080,
			"maxConnections":             8192,
			"maxHeaderCount":             100,
			"maxKeepAliveRequests":       100,
			"maxThreads":                 300,
			"maxThreadsWithExecutor":     300,
			"minSpareThreads":            10,
			"modelerType":                "org.apache.tomcat.util.net.AprEndpoint",
			"name":                       "\"http-apr-8080\"",
			"paused":                     false,
			"pollTime":                   2000,
			"port":                       8080,
			"running":                    true,
			"sSLCipherSuite":             "HIGH:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!kRSA",
			"sSLDisableCompression":      false,
			"sSLEnabled":                 false,
			"sSLHonorCipherOrder":        false,
			"sSLInsecureRenegotiation":   false,
			"sSLProtocol":                "all",
			"sSLVerifyClient":            "none",
			"sSLVerifyDepth":             10,
			"sendfileCount":              0,
			"sendfileSize":               1024,
			"sendfileThreadCount":        0,
			"sessionTimeout":             "86400",
			"soLinger":                   -1,
			"soTimeout":                  20000,
			"sslProtocol":                "TLS",
			"tcpNoDelay":                 true,
			"threadPriority":             5,
			"useComet":                   true,
			"useCometTimeout":            false,
			"usePolling":                 true,
			"useSendfile":                true,
			"useServerCipherSuitesOrder": "",
		},
		map[string]interface{}{
			"Catalina:type":          "ThreadPool",
			"acceptorThreadCount":    1,
			"acceptorThreadPriority": 5,
			"algorithm":              "SunX509",
			"backlog":                100,
			"bean":                   "type=ThreadPool,name=\"ajp-apr-8009\"",
			"bindOnInit":             true,
			"clientAuth":             "false",
			"connectionCount":        1,
			"currentThreadCount":     10,
			"currentThreadsBusy":     0,
			"daemon":                 true,
			"deferAccept":            true,
			"domain":                 "Catalina",
			"executorTerminationTimeoutMillis": 5000,
			"hi":                         "hello",
			"keepAliveCount":             0,
			"keepAliveTimeout":           -1,
			"keystoreFile":               "/root/.keystore",
			"keystoreType":               "JKS",
			"localPort":                  8009,
			"maxConnections":             8192,
			"maxHeaderCount":             100,
			"maxKeepAliveRequests":       100,
			"maxThreads":                 200,
			"maxThreadsWithExecutor":     200,
			"minSpareThreads":            10,
			"modelerType":                "org.apache.tomcat.util.net.AprEndpoint",
			"name":                       "\"ajp-apr-8009\"",
			"paused":                     false,
			"pollTime":                   2000,
			"port":                       8009,
			"running":                    true,
			"sSLCipherSuite":             "HIGH:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!MD5:!kRSA",
			"sSLDisableCompression":      false,
			"sSLEnabled":                 false,
			"sSLHonorCipherOrder":        false,
			"sSLInsecureRenegotiation":   false,
			"sSLProtocol":                "all",
			"sSLVerifyClient":            "none",
			"sSLVerifyDepth":             10,
			"sendfileCount":              0,
			"sendfileSize":               1024,
			"sendfileThreadCount":        0,
			"sessionTimeout":             "86400",
			"soLinger":                   -1,
			"soTimeout":                  -1,
			"sslProtocol":                "TLS",
			"tcpNoDelay":                 true,
			"threadPriority":             5,
			"useComet":                   true,
			"useCometTimeout":            false,
			"usePolling":                 true,
			"useSendfile":                false,
			"useServerCipherSuitesOrder": "",
		},
	}

	ParseJMX(dataInterface, config.APIs[0].Commands[0], &prevSample)

	if len(expectedDatastore) != len(load.Store.Data) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(load.Store.Data))
	}

	for i, sample := range expectedDatastore {
		switch sample := sample.(type) {
		case map[string]interface{}:
			for key := range sample {
				switch recSample := load.Store.Data[i].(type) {
				case map[string]interface{}:
					if sample["bean"] == recSample["bean"] {
						if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
							t.Errorf("%v want %v, got %v", key, sample[key], recSample[key])
						}
					}
				}
			}
		}
	}
}

func TestParseJMXwithCompressBean(t *testing.T) {
	load.Refresh()
	// create a listener with desired port
	b, _ := ioutil.ReadFile("../../test/payloads/tomcatJMX.out")
	_, dataInterface := detectCommandOutput(string(b), "jmx")
	prevSample := map[string]interface{}{
		"hi": "hello",
	}
	config := load.Config{
		Name: "jmxFlex",
		APIs: []load.API{
			{
				Name: "tomcat",
				Commands: []load.Command{
					{
						Output:       "jmx",
						CompressBean: true,
						Run:          "Catalina:type=ThreadPool,name=*",
						Jmx: load.JMX{
							Port:           "9001",
							User:           "batman",
							Pass:           "robin",
							KeyStore:       "abc",
							KeyStorePass:   "def",
							TrustStore:     "abc",
							TrustStorePass: "def",
						},
					},
				},
			},
		},
	}
	expectedDatastore := []interface{}{
		map[string]interface{}{
			`ThreadPool."ajp-apr-8009".acceptorThreadCount`:              1,
			`ThreadPool."ajp-apr-8009".acceptorThreadPriority`:           5,
			`ThreadPool."ajp-apr-8009".algorithm`:                        "SunX509",
			`ThreadPool."ajp-apr-8009".bindOnInit`:                       true,
			`ThreadPool."ajp-apr-8009".currentThreadsBusy`:               0,
			`ThreadPool."ajp-apr-8009".daemon`:                           true,
			`ThreadPool."ajp-apr-8009".executorTerminationTimeoutMillis`: 5000,
			`ThreadPool."ajp-apr-8009".keystoreFile`:                     `/root/.keystore`,
			`ThreadPool."ajp-apr-8009".keystoreType`:                     `JKS`,
			`ThreadPool."ajp-apr-8009".maxConnections`:                   8192,
			`ThreadPool."ajp-apr-8009".maxHeaderCount`:                   100,
			`ThreadPool."ajp-apr-8009".maxKeepAliveRequests`:             100,
			`ThreadPool."ajp-apr-8009".maxThreads`:                       200,
			`ThreadPool."ajp-apr-8009".maxThreadsWithExecutor`:           200,
			`ThreadPool."ajp-apr-8009".minSpareThreads`:                  10,
			`ThreadPool."ajp-apr-8009".modelerType`:                      `org.apache.tomcat.util.net.AprEndpoint`,
			`ThreadPool."ajp-apr-8009".name`:                             `ajp-apr-8009`,
			`ThreadPool."ajp-apr-8009".pollTime`:                         2000,
			`ThreadPool."ajp-apr-8009".running`:                          true,
			`ThreadPool."ajp-apr-8009".sSLEnabled`:                       false,
			`ThreadPool."ajp-apr-8009".sSLHonorCipherOrder`:              false,
			`ThreadPool."ajp-apr-8009".sSLVerifyClient`:                  `none`,
			`ThreadPool."ajp-apr-8009".sendfileSize`:                     1024,
			`ThreadPool."ajp-apr-8009".sendfileThreadCount`:              0,
			`ThreadPool."ajp-apr-8009".sessionTimeout`:                   `86400`,
			`ThreadPool."ajp-apr-8009".soLinger`:                         -1,
			`ThreadPool."ajp-apr-8009".tcpNoDelay`:                       true,
			`ThreadPool."ajp-apr-8009".useComet`:                         true,
			`ThreadPool."ajp-apr-8009".useCometTimeout`:                  false,
			`ThreadPool."ajp-apr-8009".usePolling`:                       true,
			`ThreadPool."http-apr-8080".algorithm`:                       `SunX509`,
			`ThreadPool."http-apr-8080".backlog`:                         100,
			`ThreadPool."http-apr-8080".connectionCount`:                 1,
			`ThreadPool."http-apr-8080".currentThreadCount`:              10,
			`ThreadPool."http-apr-8080".currentThreadsBusy`:              0,
			`ThreadPool."http-apr-8080".keepAliveCount`:                  0,
			`ThreadPool."http-apr-8080".keystoreType`:                    `JKS`,
			`ThreadPool."http-apr-8080".maxConnections`:                  8192,
			`ThreadPool."http-apr-8080".maxThreads`:                      300,
			`ThreadPool."http-apr-8080".minSpareThreads`:                 "10",
			`ThreadPool."http-apr-8080".modelerType`:                     `org.apache.tomcat.util.net.AprEndpoint`,
			`ThreadPool."http-apr-8080".name`:                            `http-apr-8080`,
			`ThreadPool."http-apr-8080".sSLDisableCompression`:           false,
			`ThreadPool."http-apr-8080".sSLEnabled`:                      false,
			`ThreadPool."http-apr-8080".sSLInsecureRenegotiation`:        false,
			`ThreadPool."http-apr-8080".sSLVerifyClient`:                 `none`,
			`ThreadPool."http-apr-8080".sSLVerifyDepth`:                  10,
			`ThreadPool."http-apr-8080".sendfileCount`:                   0,
			`ThreadPool."http-apr-8080".soLinger`:                        -1,
			`ThreadPool."http-apr-8080".soTimeout`:                       20000,
			`ThreadPool."http-apr-8080".tcpNoDelay`:                      true,
			`ThreadPool."http-apr-8080".threadPriority`:                  5,
			`ThreadPool."http-apr-8080".useComet`:                        true,
			`ThreadPool."http-apr-8080".useSendfile`:                     true,
		},
	}

	ParseJMX(dataInterface, config.APIs[0].Commands[0], &prevSample)

	if len(expectedDatastore) != len(load.Store.Data) {
		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(load.Store.Data))
		t.Errorf("%v", (load.Store.Data))

	}

	if len(load.Store.Data[0].(map[string]interface{})) != 102 {
		t.Errorf("Incorrect number of metrics generated expected: %d, got: %d", 102, len(load.Store.Data[0].(map[string]interface{})))
	}

	for i, sample := range expectedDatastore {
		switch sample := sample.(type) {
		case map[string]interface{}:
			for key := range sample {
				switch recSample := load.Store.Data[i].(type) {
				case map[string]interface{}:
					if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
						t.Errorf("%v want %v, got %v", key, sample[key], recSample[key])
					}
				}
			}
		}
	}

}
