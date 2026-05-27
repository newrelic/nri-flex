/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestGetHostKeyCallback_MissingKnownHostsFile(t *testing.T) {
	scp := load.SCP{KnownHostsFile: "/nonexistent/path/known_hosts"}
	_, err := getHostKeyCallback(scp)
	if err == nil {
		t.Fatal("expected error for missing known_hosts file")
	}
	if !strings.Contains(err.Error(), "known_hosts file not found") {
		t.Fatalf("expected 'known_hosts file not found' error, got: %v", err)
	}
}

func TestGetHostKeyCallback_CustomKnownHostsFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scp-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a valid (empty) known_hosts file
	knownHostsPath := filepath.Join(tmpDir, "known_hosts")
	if err := os.WriteFile(knownHostsPath, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	scp := load.SCP{KnownHostsFile: knownHostsPath}
	callback, err := getHostKeyCallback(scp)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if callback == nil {
		t.Fatal("expected non-nil callback")
	}
}

func TestGetHostKeyCallback_DefaultPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}

	defaultPath := filepath.Join(home, ".ssh", "known_hosts")
	scp := load.SCP{}

	_, cbErr := getHostKeyCallback(scp)
	if _, statErr := os.Stat(defaultPath); os.IsNotExist(statErr) {
		// If default known_hosts doesn't exist, we expect an error
		if cbErr == nil {
			t.Fatal("expected error when default known_hosts doesn't exist")
		}
		if !strings.Contains(cbErr.Error(), "known_hosts file not found") {
			t.Fatalf("expected 'known_hosts file not found' error, got: %v", cbErr)
		}
	} else {
		// If it exists, callback should succeed
		if cbErr != nil {
			t.Fatalf("expected no error when default known_hosts exists, got: %v", cbErr)
		}
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/foo/bar", filepath.Join(home, "foo/bar")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
		{"", ""},
	}

	for _, tc := range tests {
		result := expandHome(tc.input)
		if result != tc.expected {
			t.Errorf("expandHome(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

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
			"known_hosts file not found",
			"host key verification",
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
