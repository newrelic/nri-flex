/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// RunScpWithTimeout performs scp with timeout to gather data from a remote file.
func RunScpWithTimeout(dataStore *[]interface{}, cfg *load.Config, api load.API) error {
	load.Logrus.Debugf("%v - running scp requests", cfg.Name)
	remoteFile := api.Scp.RemoteFile

	client, err := getSSHConnection(cfg, api)
	if err != nil {
		return err
	}

	srcFile, err := client.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("ssh: failed to open source file: %s, error: %v", remoteFile, err)
	}

	fileContent, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("ssh: failed to read file: %s, error: %v", remoteFile, err)
	}

	return handleScpJSON(dataStore, fileContent)
}

func getSSHConnection(yml *load.Config, api load.API) (*sftp.Client, error) {
	var user string
	var timeout time.Duration

	host := api.Scp.Host
	port := api.Scp.Port

	if yml.Global.User != "" {
		user = yml.Global.User
	}

	if api.Scp.User != "" {
		user = api.Scp.User
	}

	if yml.Global.Timeout > 0 {
		timeout = time.Duration(yml.Global.Timeout) * time.Millisecond
	} else {
		timeout = load.DefaultPingTimeout
	}

	authMethod, err := getAuthMethod(yml, api)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
	} // #nosec

	sshConfig.SetDefaults()

	conn, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh: failed to connect to sftp host: %s, with user %s, error: %v",
			host, user, err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("ssh: failed to init sftp client, error: %v", err)
	}
	return client, nil
}

func getAuthMethod(yml *load.Config, api load.API) (ssh.AuthMethod, error) {
	var sshPemFile, pass, passphrase string

	if yml.Global.SSHPEMFile != "" {
		sshPemFile = yml.Global.SSHPEMFile
	}
	if api.Scp.SSHPEMFile != "" {
		sshPemFile = api.Scp.SSHPEMFile
	}

	if sshPemFile != "" {
		return publicKeyFile(sshPemFile)
	}

	if yml.Global.Pass != "" {
		pass = yml.Global.Pass
	}
	if yml.Global.Passphrase != "" {
		passphrase = yml.Global.Passphrase
	}

	if api.Scp.Pass != "" {
		pass = api.Scp.Pass
	}
	if api.Scp.Passphrase != "" {
		passphrase = api.Scp.Passphrase
	}

	if passphrase != "" {
		encryptedPass, err := hex.DecodeString(pass)
		if err == nil {
			realPass, err := utils.Decrypt(encryptedPass, passphrase)
			if err == nil {
				pass = string(realPass)
			}
		}
	}

	return ssh.Password(pass), nil
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read ssh pem file: %v", err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return ssh.PublicKeys(key), nil
}

func handleScpJSON(dataStore *[]interface{}, body []byte) error {
	newBody := strings.Replace(string(body), " ", "", -1)
	var data interface{}
	err := json.Unmarshal([]byte(newBody), &data)
	if err != nil {
		return fmt.Errorf("ssh: failed to unmarshal JSON error: %v", err)
	}
	*dataStore = append(*dataStore, data)
	return nil
}
