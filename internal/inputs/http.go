package inputs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/parnurzeal/gorequest"
)

// RunHTTP Executes HTTP Requests
func RunHTTP(dataStore *[]interface{}, doLoop *bool, yml *load.Config, api load.API, reqURL *string) {
	logger.Flex("debug", nil, fmt.Sprintf("%v - running http requests", yml.Name), false)
	for *doLoop {
		request := gorequest.New()

		if api.EscapeURL {
			*reqURL = url.QueryEscape(*reqURL)
		}

		// weird edge case, happens with rabbitmq
		if strings.HasSuffix(*reqURL, "//") {
			*reqURL = strings.TrimSuffix(*reqURL, "//")
			*reqURL += "/%2f"
		}

		*reqURL = yml.Global.BaseURL + *reqURL
		switch {
		case api.Method == "POST" && api.Payload != "":
			request = request.Post(*reqURL)
			request = request.Send(api.Payload)
		case api.Method == "PUT" && api.Payload != "":
			request = request.Put(*reqURL)
			request = request.Send(api.Payload)
		default:
			request = request.Get(*reqURL)
		}

		request = setRequestOptions(request, *yml, api)
		logger.Flex("debug", nil, fmt.Sprintf("sending %v request to %v", request.Method, *reqURL), false)
		resp, _, errors := request.End()
		if resp != nil {
			nextLink := ""
			if resp.Header["Link"] != nil {
				headerLinks := strings.Split(resp.Header["Link"][0], ",")
				for _, link := range headerLinks {
					if strings.Contains(link, "next") {
						theLink := strings.Split(link, ";")
						nextLink = strings.Replace((strings.Replace(theLink[0], "<", "", -1)), ">", "", -1)
						nextLink = strings.TrimPrefix(nextLink, " ")
					}
				}
			}

			// responseReceived := map[string]interface{}{}
			contentType := resp.Header.Get("Content-Type")
			responseError := ""

			switch {
			case api.Prometheus.Enable:
				Prometheus(dataStore, resp.Body, yml, &api)
			case strings.Contains(contentType, "application/json"):
				body, _ := ioutil.ReadAll(resp.Body)
				handleJSON(dataStore, body, &resp, doLoop, reqURL, nextLink)
			default:
				// some apis do not specify a content-type header, if not set attempt to detect if the payload is json
				body, _ := ioutil.ReadAll(resp.Body)
				strBody := string(body)
				output, _ := detectCommandOutput(strBody, "")
				switch output {
				case load.TypeJSON:
					handleJSON(dataStore, body, &resp, doLoop, reqURL, nextLink)
				default:
					logger.Flex("debug", fmt.Errorf("%v - Not sure how to handle this payload? ContentType: %v", api.URL, contentType), "", false)
					logger.Flex("debug", fmt.Errorf("%v - storing unknown http output into datastore", api.URL), "", false)
					if yml.Datastore == nil {
						yml.Datastore = map[string][]interface{}{}
					}
					yml.Datastore[api.URL] = []interface{}{
						map[string]interface{}{
							"http": strBody,
						},
					}
				}
			}

			if responseError == "" {
				if nextLink != "" {
					*reqURL = nextLink
				} else {
					*doLoop = false
				}
			}
		} else {
			for _, err := range errors {
				logger.Flex("debug", err, "", false)
			}
			*doLoop = false
		}
	}
}

// setRequestOptions
// Sets global config for all APIs/Endpoints
// However, nested configs that are defined will take precedence over global config
func setRequestOptions(request *gorequest.SuperAgent, yml load.Config, api load.API) *gorequest.SuperAgent {
	rootCAs := x509.NewCertPool()
	if yml.Global.Timeout > 0 {
		request = request.Timeout(time.Duration(yml.Global.Timeout) * time.Millisecond)
	}
	if yml.Global.Proxy != "" {
		request = request.Proxy(yml.Global.Proxy)
	}
	if yml.Global.User != "" {
		request = request.SetBasicAuth(yml.Global.User, yml.Global.Pass)
	}
	if yml.Global.TLSConfig.Ca != "" {
		ca, err := ioutil.ReadFile(yml.Global.TLSConfig.Ca)
		if err != nil {
			logger.Flex("error", err, "failed to read ca", false)
		} else {
			rootCAs.AppendCertsFromPEM(ca)
		}
	}
	for h, v := range yml.Global.Headers {
		request = request.Set(h, v)
	}
	if api.Timeout > 0 {
		request = request.Timeout(time.Duration(api.Timeout) * time.Millisecond)
	}
	if api.Proxy != "" {
		request = request.Proxy(api.Proxy)
	}
	if api.User != "" {
		request = request.SetBasicAuth(api.User, api.Pass)
	}
	for h, v := range api.Headers {
		request = request.Set(h, v)
	}
	if api.TLSConfig.Ca != "" {
		ca, err := ioutil.ReadFile(api.TLSConfig.Ca)
		if err != nil {
			logger.Flex("error", err, "failed to read ca", false)
		} else {
			rootCAs.AppendCertsFromPEM(ca)
		}
	}

	request = request.TLSClientConfig(&tls.Config{
		InsecureSkipVerify: yml.Global.TLSConfig.InsecureSkipVerify,
		MinVersion:         yml.Global.TLSConfig.MinVersion,
		MaxVersion:         yml.Global.TLSConfig.MaxVersion,
		RootCAs:            rootCAs,
	})

	if api.TLSConfig.Enable {
		request = request.TLSClientConfig(&tls.Config{
			InsecureSkipVerify: api.TLSConfig.InsecureSkipVerify,
			MinVersion:         api.TLSConfig.MinVersion,
			MaxVersion:         api.TLSConfig.MaxVersion,
			RootCAs:            rootCAs,
		})
	}

	return request
}

// handleJSON Process JSON Payload
func handleJSON(dataStore *[]interface{}, body []byte, resp *gorequest.Response, doLoop *bool, url *string, nextLink string) {
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		logger.Flex("error", err, "", false)
	} else {
		switch f := f.(type) {
		case []interface{}:
			for _, sample := range f {

				switch sample := sample.(type) {
				case map[string]interface{}:
					httpSample := sample
					httpSample["api.StatusCode"] = (*resp).StatusCode
					*dataStore = append(*dataStore, httpSample)
					// load.StoreAppend(httpSample)
				case string:
					strSample := map[string]interface{}{
						"output": sample,
					}
					// load.StoreAppend(strSample)
					*dataStore = append(*dataStore, strSample)
				default:
					logger.Flex("debug", fmt.Errorf("not sure how to handle this %v", sample), "", false)
				}
			}

		case map[string]interface{}:
			theSample := f
			theSample["api.StatusCode"] = (*resp).StatusCode
			// load.StoreAppend(theSample)
			*dataStore = append(*dataStore, theSample)

			if theSample["error"] != nil {
				logger.Flex("debug", nil, "Request failed "+fmt.Sprintf("%v", theSample["error"]), false)
			}

			if theSample["error"] == nil && nextLink != "" {
				*url = nextLink
			} else {
				*doLoop = false
			}
		}
	}
}
