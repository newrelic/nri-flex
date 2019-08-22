package inputs

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
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

		handlePagination(reqURL, &api.Pagination, nil, nil, 200)
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

			contentType := resp.Header.Get("Content-Type")
			responseError := ""

			logger.Flex("debug", nil, fmt.Sprintf("URL: %v Status: %v Code: %d", *reqURL, resp.Status, resp.StatusCode), false)

			switch {
			case api.Prometheus.Enable:
				Prometheus(dataStore, resp.Body, yml, &api)
			case strings.Contains(contentType, "application/json"):
				body, _ := ioutil.ReadAll(resp.Body)
				addPage := handlePagination(nil, &api.Pagination, &nextLink, body, resp.StatusCode)
				if api.Debug {
					logger.Flex("debug", nil, fmt.Sprintf("HTTP Debug:\nURL: %v\nBody:\n%v\n", *reqURL, string(body)), false)
				}
				// if not using pagination handle json for any response, if using pagination check the status code before storing
				if api.Pagination.OriginalURL == "" || (api.Pagination.OriginalURL != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299) && addPage {
					handleJSON(dataStore, body, &resp, doLoop, reqURL, nextLink)
				}
			default:
				// some apis do not specify a content-type header, if not set attempt to detect if the payload is json
				body, err := ioutil.ReadAll(resp.Body)
				addPage := handlePagination(nil, &api.Pagination, &nextLink, body, resp.StatusCode)

				if err != nil {
					logger.Flex("error", err, fmt.Sprintf("HTTP URL: %v failed to read resp.Body", *reqURL), false)
				} else {
					strBody := string(body)
					if api.Debug {
						logger.Flex("debug", nil, fmt.Sprintf("HTTP Debug:\nURL: %v\nBody:\n%v\n", *reqURL, strBody), false)
					}
					output, _ := detectCommandOutput(strBody, "")
					switch output {
					case load.TypeJSON:
						// if not using pagination handle json for any response, if using pagination check the status code before storing
						if api.Pagination.OriginalURL == "" || (api.Pagination.OriginalURL != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299) && addPage {
							handleJSON(dataStore, body, &resp, doLoop, reqURL, nextLink)
						}
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

	tmpGlobalTLSConfig := tls.Config{
		InsecureSkipVerify: yml.Global.TLSConfig.InsecureSkipVerify,
		MinVersion:         yml.Global.TLSConfig.MinVersion,
		MaxVersion:         yml.Global.TLSConfig.MaxVersion,
	}

	if yml.Global.TLSConfig.Ca != "" {
		ca, err := ioutil.ReadFile(yml.Global.TLSConfig.Ca)
		if err != nil {
			logger.Flex("error", err, "failed to read ca", false)
		} else {
			rootCAs.AppendCertsFromPEM(ca)
			tmpGlobalTLSConfig.RootCAs = rootCAs
		}
	}

	request = request.TLSClientConfig(&tmpGlobalTLSConfig)

	if api.TLSConfig.Enable {
		tmpAPITLSConfig := tls.Config{
			InsecureSkipVerify: api.TLSConfig.InsecureSkipVerify,
			MinVersion:         api.TLSConfig.MinVersion,
			MaxVersion:         api.TLSConfig.MaxVersion,
		}

		if api.TLSConfig.Ca != "" {
			ca, err := ioutil.ReadFile(api.TLSConfig.Ca)
			if err != nil {
				logger.Flex("error", err, "failed to read ca", false)
			} else {
				rootCAs.AppendCertsFromPEM(ca)
				tmpAPITLSConfig.RootCAs = rootCAs
			}
		}
		request = request.TLSClientConfig(&tmpAPITLSConfig)
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

func handlePagination(url *string, Pagination *load.Pagination, nextLink *string, body []byte, code int) bool {
	if url != nil && strings.Contains(*url, "${page}") && (code >= 200 && code <= 299) {
		(*Pagination).OriginalURL = *url
		(*Pagination).NoPages = 1
		(*Pagination).PageMarker = Pagination.PageStart
		if (*Pagination).Increment == 0 {
			(*Pagination).Increment = 1
		}
		*url = strings.Replace(*url, "${page}", fmt.Sprintf("%d", Pagination.PageStart), -1)
		*url = strings.Replace(*url, "${limit}", fmt.Sprintf("%d", Pagination.PageLimit), -1)
		logger.Flex("debug", nil, fmt.Sprintf("URL: %v begin pagination handling", *url), false)
	} else if Pagination.OriginalURL != "" && nextLink != nil && (code >= 200 && code <= 299) {
		if Pagination.MaxPages == 0 && Pagination.PageLimitKey == "" && Pagination.PayloadKey == "" {
			link := ""
			if url != nil {
				link = *url
			}
			if nextLink != nil {
				link = *nextLink
			}
			logger.Flex("debug", nil, fmt.Sprintf("URL: %v not walking next link, max_pages and/or payload_key, and/or page_limit_key has not been set", link), false)
		} else {
			continueRequest := true
			customPageMarker := false
			nextCursor := ""
			payloadEmpty := false
			payloadKeyFound := false
			manualNextLink := ""
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, body); err != nil {
				logger.Flex("error", err, "", false)
			} else {
				if Pagination.PageLimitKey != "" || Pagination.PageNextKey != "" || Pagination.PayloadKey != "" || Pagination.MaxPagesKey != "" || Pagination.NextCursorKey != "" {
					jsonString := buffer.String()
					if Pagination.PageLimitKey != "" { // offset
						matches := paginationRegex(fmt.Sprintf(`"%v":(\d+)|"%v":"(\d+)"`, Pagination.PageLimitKey, Pagination.PageLimitKey), jsonString, nextLink)
						if len(matches) >= 2 {
							no, nerr := strconv.Atoi(matches[1])
							if nerr != nil {
								logger.Flex("error", nerr, nil, false)
							} else {
								Pagination.PageLimit = no
							}
						}
					}
					if Pagination.MaxPagesKey != "" {
						matches := paginationRegex(fmt.Sprintf(`"%v":(\d+)|"%v":"(\d+)"`, Pagination.MaxPagesKey, Pagination.MaxPagesKey), jsonString, nextLink)
						if len(matches) >= 2 {
							no, nerr := strconv.Atoi(matches[1])
							if nerr != nil {
								logger.Flex("error", nerr, nil, false)
							} else {
								Pagination.MaxPages = no
							}
						}
					}
					if Pagination.PageNextKey != "" {
						matches := paginationRegex(fmt.Sprintf(`"%v":(\d+)|"%v":"(\d+)"`, Pagination.PageNextKey, Pagination.PageNextKey), jsonString, nextLink)
						if len(matches) >= 2 {
							no, nerr := strconv.Atoi(matches[1])
							if nerr != nil {
								logger.Flex("error", nerr, nil, false)
							} else {
								Pagination.PageMarker = no
								customPageMarker = true
							}
						}
					}
					if Pagination.NextCursorKey != "" {
						matches := paginationRegex(fmt.Sprintf(`"%v":(\d+)|"%v":"(\d+)"`, Pagination.NextCursorKey, Pagination.NextCursorKey), jsonString, nextLink)
						if len(matches) >= 2 {
							nextCursor = matches[1]
						}
					}
					if Pagination.NextLinkKey != "" {
						matches := paginationRegex(fmt.Sprintf(`"%v":\"(\S+)\"`, Pagination.NextLinkKey), jsonString, nextLink)
						if len(matches) >= 2 {
							manualNextLink = matches[1]
						}
					}
					if Pagination.PayloadKey != "" {
						matches := paginationRegex(fmt.Sprintf(`"%v":(\[(.*?)\]|\{(.*?)\})`, Pagination.PayloadKey), jsonString, nextLink)
						if len(matches) >= 3 {
							payloadKeyFound = true
							if matches[1] == "{}" || matches[1] == "[]" {
								*nextLink = ""
								continueRequest = false
								payloadEmpty = true
								logger.Flex("debug", nil, fmt.Sprintf("URL: %v walk payload %v %v empty", *nextLink, Pagination.PayloadKey, matches[1]), false)
							}
						}
					}
				}
			}

			if (Pagination.PageMarker >= Pagination.MaxPages && Pagination.PayloadKey == "" && payloadKeyFound) || (Pagination.PayloadKey != "" && payloadKeyFound && payloadEmpty) {
				logger.Flex("debug", nil, fmt.Sprintf("URL: %v max pages reached %d or payload empty %v", *nextLink, Pagination.MaxPages, payloadEmpty), false)
				*nextLink = ""
				continueRequest = false
				return false
			}
			if continueRequest {
				page := ""
				if !customPageMarker {
					(*Pagination).PageMarker = (*Pagination).PageMarker + (*Pagination).Increment
					page = fmt.Sprintf("%d", (*Pagination).PageMarker)
				}
				if nextCursor != "" {
					page = nextCursor
				}
				if page != "" && Pagination.NextLinkKey == "" {
					*nextLink = strings.Replace((*Pagination).OriginalURL, "${page}", page, -1)
					*nextLink = strings.Replace(*nextLink, "${limit}", fmt.Sprintf("%d", Pagination.PageLimit), -1)
					logger.Flex("debug", nil, fmt.Sprintf("URL: %v walking next link", *nextLink), false)
				}
				if manualNextLink != "" {
					*nextLink = manualNextLink
					logger.Flex("debug", nil, fmt.Sprintf("URL: %v walking next link", *nextLink), false)
				}
			}
		}
	}
	return true
}

// paginationRegex
func paginationRegex(regexKey string, jsonString string, nextLink *string) []string {
	re, err := regexp.Compile(regexKey)
	if err != nil {
		logger.Flex("error", err, fmt.Sprintf("URL: %v regex compile failed %v", *nextLink, regexKey), false)
	} else {
		return re.FindStringSubmatch(jsonString)
	}
	return []string{}
}
