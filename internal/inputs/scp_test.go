/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"strings"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestPublicKeyFile(t *testing.T) {
	load.Refresh()

	config := load.Config{
		Name: "scpexample",
		Global: load.Global{
			Timeout:    3000,
			SSHPEMFile: "../../test/payloads/test_key.pem",
		},
		APIs: []load.API{
			{
				Name: "scp1",
				Scp: load.SCP{
					User:       "newrelic",
					Pass:       "N3wr3lic!",
					Host:       "8.8.8.8",
					Port:       "22",
					RemoteFile: "/cmd/meta.json",
				},
			},
		},
	}

	// The error can be either RSA parsing error or SSH certificate error
	// depending on the key format and parsing library behavior
	expectedErrors := []string{
		"ssh: unsupported key type \"CERTIFICATE\"",
		"failed to parse private key: crypto/rsa: invalid CRT coefficient",
		"crypto/rsa: invalid CRT coefficient",
	}

	_, err := publicKeyFile(config.Global.SSHPEMFile)
	if err != nil {
		errorMatched := false
		actualError := err.Error()

		// Check if the actual error contains any of the expected error messages
		for _, expectedErr := range expectedErrors {
			if strings.Contains(actualError, expectedErr) {
				errorMatched = true
				break
			}
		}

		if !errorMatched {
			t.Errorf("received error '%v' does not match any expected errors: %v", err, expectedErrors)
		}
	} else {
		t.Error("expected an error but got none")
	}
}

func TestGetSSHConnection(t *testing.T) {
	load.Refresh()

	config := load.Config{
		Name: "scpexample",
		Global: load.Global{
			Timeout: 3000,
		},
		APIs: []load.API{
			{
				Name: "scp1",
				Scp: load.SCP{
					User:       "newrelic",
					Pass:       "N3wr3lic!",
					Host:       "8.8.8.8",
					Port:       "22",
					RemoteFile: "/cmd/meta.json",
				},
			},
		},
	}

	_, err := getSSHConnection(&config, config.APIs[0])

	if err != nil {
		actualError := err.Error()
		// Check for common connection error patterns since exact error messages can vary
		expectedPatterns := []string{
			"i/o timeout",
			"connection refused",
			"no route to host",
			"network is unreachable",
			"failed to connect to sftp host",
		}

		errorMatched := false
		for _, pattern := range expectedPatterns {
			if strings.Contains(actualError, pattern) {
				errorMatched = true
				break
			}
		}

		if !errorMatched {
			t.Errorf("received error '%v' does not contain expected connection error patterns: %v", err, expectedPatterns)
		}
	} else {
		t.Error("expected a connection error but got none")
	}
}
