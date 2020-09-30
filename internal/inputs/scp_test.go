/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
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
	// expectedErr := "ssh: cannot decode encrypted private keys"
	expectedErr := "ssh: unsupported key type \"CERTIFICATE\""
	_, err := publicKeyFile(config.Global.SSHPEMFile)
	if err != nil {
		if err.Error() != expectedErr {
			t.Errorf("received error  %v does not match expected %v", err, expectedErr)
		}
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
	expectedErr := "ssh: failed to connect to sftp host: 8.8.8.8, with user newrelic, error: dial tcp 8.8.8.8:22: i/o timeout"
	if err != nil {
		if err.Error() != expectedErr {
			t.Errorf("received error '%v' does not match expected '%v'", err, expectedErr)
		}
	}

}
