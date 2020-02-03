/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package outputs

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
)

// postRequest wraps request and attaches needed headers and zlib compression
func postRequest(url string, key string, data []byte) error {
	var zlibCompressedPayload bytes.Buffer
	w := zlib.NewWriter(&zlibCompressedPayload)
	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("http: failed to compress payload, %v", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("http: failed to close zlib writer, %v", err)
	}

	load.Logrus.Debugf("http: insights - bytes %d events %d", len(zlibCompressedPayload.Bytes()), len(load.Entity.Metrics))

	tr := &http.Transport{IdleConnTimeout: 15 * time.Second}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(zlibCompressedPayload.Bytes()))
	if err != nil {
		return fmt.Errorf("http: unable to create http.Request, %v", err)
	}

	req.Header.Set("Content-Encoding", "deflate")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", key)

	resp, err := client.Do(req)
	if err != nil {
		load.Logrus.WithError(err).Error("http: failed to send")
	}
	defer resp.Body.Close()

	if resp == nil {
		return errors.New("http: response nil")
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		err = fmt.Errorf("http: post failed, status code: %d", resp.StatusCode)
	}
	return err
}
