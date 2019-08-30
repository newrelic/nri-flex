/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"net/http"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// postRequest wraps request and attaches needed headers and zlib compression
func postRequest(url string, key string, data []byte) {
	var zlibCompressedPayload bytes.Buffer
	w := zlib.NewWriter(&zlibCompressedPayload)
	_, err := w.Write(data)
	logger.Flex("error", err, "unable to write zlib compressed form", false)
	logger.Flex("error", w.Close(), "unable to close zlib writer", false)
	if err != nil {
		logger.Flex("error", fmt.Errorf("failed to compress payload"), "", false)
	} else {
		tr := &http.Transport{IdleConnTimeout: 15 * time.Second}
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(zlibCompressedPayload.Bytes()))
		logger.Flex("debug", nil, fmt.Sprintf("insights: bytes %d events %d", len(zlibCompressedPayload.Bytes()), len(load.Entity.Metrics)), false)

		if err != nil {
			logger.Flex("error", err, "unable to create http.Request", false)
		} else {
			req.Header.Set("Content-Encoding", "deflate")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Insert-Key", key)
			resp, err := client.Do(req)
			logger.Flex("error", err, "unable to send", false)
			if resp != nil {
				if resp.StatusCode > 299 || resp.StatusCode < 200 {
					logger.Flex("error", fmt.Errorf("http post unsuccessful code %d", resp.StatusCode), "", false)
				}
			} else {
				logger.Flex("error", fmt.Errorf("http response nil"), "", false)
			}
		}
	}
}
