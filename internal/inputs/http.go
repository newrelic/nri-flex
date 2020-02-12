/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

// RunHTTP Executes HTTP Requests
func RunHTTP(dataStore *[]interface{}, doLoop *bool, yml *load.Config, api load.API, reqURL *string) {
	load.Logrus.Debugf("%v - running http requests", yml.Name)
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
		case api.Method == http.MethodPost && api.Payload != "":
			request = request.Post(*reqURL)
			request = request.Send(api.Payload)
		case api.Method == http.MethodPut && api.Payload != "":
			request = request.Put(*reqURL)
			request = request.Send(api.Payload)
		default:
			request = request.Get(*reqURL)
		}

		request = setRequestOptions(request, *yml, api)
		load.Logrus.Debugf("sending %v request to %v", request.Method, *reqURL)
		resp, _, errors := request.End()
		load.StatusCounterIncrement("HttpRequests")
		if resp != nil {
			nextLink := ""
			if resp.Header["Link"] != nil {
				headerLinks := strings.Split(resp.Header["Link"][0], ",")
				for _, link := range headerLinks {
					if strings.Contains(link, "next") {
						theLink := strings.Split(link, ";")
						nextLink = strings.Replace(strings.Replace(theLink[0], "<", "", -1), ">", "", -1)
						nextLink = strings.TrimPrefix(nextLink, " ")
					}
				}
			}

			contentType := resp.Header.Get("Content-Type")
			responseError := ""

			load.Logrus.Debugf("URL: %v Status: %v Code: %d", *reqURL, resp.Status, resp.StatusCode)

			switch {
			case api.Prometheus.Enable:
				Prometheus(dataStore, resp.Body, yml, &api)
			case contentType == "application/json":
				body, _ := ioutil.ReadAll(resp.Body)
				addPage := handlePagination(nil, &api.Pagination, &nextLink, body, resp.StatusCode)
				if api.Debug {
					load.Logrus.Debugf("HTTP Debug:\nURL: %v\nBody:\n%v\n", *reqURL, string(body))
				}
				// if not using pagination handle json for any response, if using pagination check the status code before storing
				if api.Pagination.OriginalURL == "" || (api.Pagination.OriginalURL != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299) && addPage {
					handleJSON(dataStore, body, &resp, doLoop, reqURL, nextLink)
				}
			case contentType == "text/xml" || contentType == "application/xml":
				jsonBody, err := xj.Convert(resp.Body)
				if err != nil {
					load.Logrus.WithError(err).Errorf("http: URL %v failed to convert XML to Json resp.Body", *reqURL)
				} else {
					if api.Pagination.OriginalURL == "" || (api.Pagination.OriginalURL != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299) {
						handleJSON(dataStore, jsonBody.Bytes(), &resp, doLoop, reqURL, nextLink)
					}
				}
			default:
				// some apis do not specify a content-type header, if not set attempt to detect if the payload is json
				body, err := ioutil.ReadAll(resp.Body)
				addPage := handlePagination(nil, &api.Pagination, &nextLink, body, resp.StatusCode)

				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"err": err,
					}).Errorf("http: URL %v failed to read resp.Body", *reqURL)
				} else {
					strBody := string(body)
					if api.Debug {
						load.Logrus.Debugf("HTTP Debug:\nURL: %v\nBody:\n%v\n", *reqURL, strBody)
					}
					output, _ := detectCommandOutput(strBody, "")
					switch output {
					case load.TypeJSON:
						// if not using pagination handle json for any response, if using pagination check the status code before storing
						if api.Pagination.OriginalURL == "" || (api.Pagination.OriginalURL != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299) && addPage {
							handleJSON(dataStore, body, &resp, doLoop, reqURL, nextLink)
						}
						// if it is XML, convert XML to JSON and process it
					case load.TypeXML:
						xmlBody := strings.NewReader(strBody)
						jsonBody, err := xj.Convert(xmlBody)

						if err != nil {
							load.Logrus.WithFields(logrus.Fields{
								"err": err,
							}).Errorf("http: URL %v failed to convert XML to Json resp.Body", *reqURL)
						} else {
							if api.Pagination.OriginalURL == "" || (api.Pagination.OriginalURL != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299) && addPage {
								handleJSON(dataStore, jsonBody.Bytes(), &resp, doLoop, reqURL, nextLink)
							}
						}
					default:
						load.Logrus.Debugf("%v - unsupported payload format: ContentType: %v", api.URL, contentType)
						load.Logrus.Debugf("%v - storing unknown http output into datastore", api.URL)

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
			httpErrorSample := map[string]interface{}{}

			for i, err := range errors {
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Debug("http: error")

				if i == 0 {
					httpErrorSample["error"] = err
				} else {
					httpErrorSample[fmt.Sprintf("error.%d", i)] = err
				}
			}

			*dataStore = append(*dataStore, httpErrorSample)
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
			load.Logrus.WithError(err).Error("http: failed to read ca")
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
				load.Logrus.WithError(err).Error("http: failed to read ca")
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
		load.Logrus.WithError(err).Error("http: failed to unmarshal json")
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
					load.Logrus.Debugf("http: unsupported sample type: %T %v", sample, sample)
				}
			}

		case map[string]interface{}:
			theSample := f
			theSample["api.StatusCode"] = (*resp).StatusCode
			// load.StoreAppend(theSample)
			*dataStore = append(*dataStore, theSample)

			if theSample["error"] != nil && fmt.Sprintf("%v", theSample["error"]) != "false" {
				load.Logrus.Debugf("http: request failed %v", theSample["error"])
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
		load.Logrus.Debugf("URL: %v begin pagination handling", *url)
	} else if Pagination.OriginalURL != "" && nextLink != nil && (code >= 200 && code <= 299) {
		if Pagination.MaxPages == 0 && Pagination.PageLimitKey == "" && Pagination.PayloadKey == "" {
			link := ""
			if url != nil {
				link = *url
			}
			load.Logrus.Debugf("URL: %v not walking next link, max_pages and/or payload_key, and/or page_limit_key has not been set", link)
		} else {
			continueRequest := true
			customPageMarker := false
			nextCursor := ""
			payloadEmpty := false
			payloadKeyFound := false
			manualNextLink := ""
			buffer := new(bytes.Buffer)
			if err := json.Compact(buffer, body); err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("http: failed to compact json")
			} else {
				if Pagination.PageLimitKey != "" || Pagination.PageNextKey != "" || Pagination.PayloadKey != "" || Pagination.MaxPagesKey != "" || Pagination.NextCursorKey != "" {
					jsonString := buffer.String()
					if Pagination.PageLimitKey != "" { // offset
						matches := paginationRegex(fmt.Sprintf(`"%v":(\d+)|"%v":"(\d+)"`, Pagination.PageLimitKey, Pagination.PageLimitKey), jsonString, nextLink)
						if len(matches) >= 2 {
							no, nerr := strconv.Atoi(matches[1])
							if nerr != nil {
								load.Logrus.WithError(nerr).Error("http: pagination failed to convert to int")
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
								load.Logrus.WithError(nerr).Error("http: pagination failed to convert to int")
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
								load.Logrus.WithError(nerr).Error("http: pagination failed to convert to int")
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
								load.Logrus.Debugf("URL: %v walk payload %v %v empty", *nextLink, Pagination.PayloadKey, matches[1])
							}
						}
					}
				}
			}

			if (Pagination.PageMarker >= Pagination.MaxPages && Pagination.PayloadKey == "" && payloadKeyFound) || (Pagination.PayloadKey != "" && payloadKeyFound && payloadEmpty) {
				load.Logrus.Debugf("URL: %v max pages reached %d or payload empty %v", *nextLink, Pagination.MaxPages, payloadEmpty)
				*nextLink = ""
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
					load.Logrus.Debugf("URL: %v walking next link", *nextLink)
				}
				if manualNextLink != "" {
					*nextLink = manualNextLink
					load.Logrus.Debugf("URL: %v walking next link", *nextLink)
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
		load.Logrus.WithError(err).Errorf("http: URL %v regex compile failed %v", *nextLink, regexKey)
	} else {
		return re.FindStringSubmatch(jsonString)
	}
	return []string{}
}
