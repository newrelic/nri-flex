package parser

import (
	"fmt"
	"io/ioutil"
	"nri-flex/internal/load"
	"nri-flex/internal/logger"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// RunHTTP Executes HTTP Requests
func RunHTTP(doLoop *bool, yml *load.Config, api load.API, reqURL *string, dataStore *[]interface{}) {
	for *doLoop {
		request := gorequest.New()
		request = request.Get(yml.Global.BaseURL + *reqURL)
		if api.Method == "POST" {
			request = request.Post(yml.Global.BaseURL + *reqURL)
			if api.Payload != "" {
				request = request.Send(api.Payload)
			}
		}
		request = setRequestOptions(request, *yml, api)

		resp, _, errors := request.End()
		if resp != nil {
			nextLink := ""
			if resp.Header["Link"] != nil {
				headerLinks := strings.Split(resp.Header["Link"][0], ",")
				for _, link := range headerLinks {
					if strings.Contains(link, "next") {
						theLink := strings.Split(link, ";")
						nextLink = strings.Replace((strings.Replace(theLink[0], "<", "", -1)), ">", "", -1)
					}
				}
			}

			// responseReceived := map[string]interface{}{}
			contentType := resp.Header.Get("Content-Type")
			responseError := ""

			switch {
			case api.Prometheus.Enable:
				Prometheus(resp.Body, dataStore, &api)
			case strings.Contains(contentType, "application/json"):
				body, _ := ioutil.ReadAll(resp.Body)
				handleJSON(body, dataStore, &resp, doLoop, reqURL, nextLink)
			default:
				// consider adding support for non JSON
				logger.Flex("debug", fmt.Errorf("%v - Not sure how to handle this payload? ContentType: %v", api.URL, contentType), "", false)
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
	if yml.Global.Timeout > 0 {
		request = request.Timeout(time.Duration(yml.Global.Timeout) * time.Millisecond)
	}
	if yml.Global.Proxy != "" {
		request = request.Proxy(yml.Global.Proxy)
	}
	if yml.Global.User != "" {
		request = request.SetBasicAuth(yml.Global.User, yml.Global.Pass)
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
	return request
}

// handleJSON Process JSON Payload
func handleJSON(body []byte, dataStore *[]interface{}, resp *gorequest.Response, doLoop *bool, url *string, nextLink string) {
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		logger.Flex("debug", err, "", false)
	} else {
		switch f := f.(type) {
		case []interface{}:
			for _, sample := range f {
				theSample := sample.(map[string]interface{})
				theSample["api.StatusCode"] = (*resp).StatusCode
				*dataStore = append(*dataStore, theSample)
			}

		case map[string]interface{}:
			theSample := f
			theSample["api.StatusCode"] = (*resp).StatusCode
			*dataStore = append(*dataStore, theSample)
			// requestedPage = f.(map[string]interface{})
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
