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
	if err != nil {
		logger.Flex("error", fmt.Errorf("failed to compress payload"), "", false)
		return
	}
	err = w.Close()
	if err != nil {
		logger.Flex("error", err, "unable to close zlib writer", false)
		return
	}
	logger.Flex("debug", nil, fmt.Sprintf("insights: bytes %d events %d", len(zlibCompressedPayload.Bytes()), len(load.Entity.Metrics)), false)

	tr := &http.Transport{IdleConnTimeout: 15 * time.Second}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(zlibCompressedPayload.Bytes()))
	if err != nil {
		logger.Flex("error", err, "unable to create http.Request", false)
		return
	}

	req.Header.Set("Content-Encoding", "deflate")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", key)

	resp, err := client.Do(req)
	if err != nil {
		logger.Flex("error", err, "unable to send", false)
	}
	defer resp.Body.Close()

	if resp != nil {
		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			logger.Flex("error", fmt.Errorf("http post unsuccessful code %d", resp.StatusCode), "", false)
		}
	} else {
		logger.Flex("error", fmt.Errorf("http response nil"), "", false)
	}

}
