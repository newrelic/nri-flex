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
	"github.com/sirupsen/logrus"
)

// postRequest wraps request and attaches needed headers and zlib compression
func postRequest(url string, key string, data []byte) {
	var zlibCompressedPayload bytes.Buffer
	w := zlib.NewWriter(&zlibCompressedPayload)
	_, err := w.Write(data)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("http: failed to compress payload")
		return
	}
	err = w.Close()
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("http: failed to close zlib writer")
		return
	}

	load.Logrus.Debug(fmt.Sprintf("http: insights - bytes %d events %d", len(zlibCompressedPayload.Bytes()), len(load.Entity.Metrics)))

	tr := &http.Transport{IdleConnTimeout: 15 * time.Second}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(zlibCompressedPayload.Bytes()))
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("http: unable to create http.Request")
		return
	}

	req.Header.Set("Content-Encoding", "deflate")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", key)

	resp, err := client.Do(req)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("http: failed to send")
	}
	defer resp.Body.Close()

	if resp != nil {
		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			load.Logrus.WithFields(logrus.Fields{
				"code": resp.StatusCode,
			}).Error("http: post failed")
		}
	} else {
		load.Logrus.Error("http: response nil")
	}

}
