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
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// RunScpWithTimeout performs scp  with timeout
func RunScpWithTimeout(dataStore *[]interface{}, cfg *load.Config, api load.API, apiNo int) {
	load.Logrus.Debug(fmt.Sprintf("%v - running scp requests", cfg.Name))
	remotefile := api.Scp.RemoteFile
	client, err := getSSHConnection(cfg, api)

	if err == nil {
		srcFile, err := client.Open(remotefile)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err":  err,
				"file": remotefile,
			}).Error("ssh: failed to open source file")
		} else {

			fileData, err := ioutil.ReadAll(srcFile)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": cfg.Name,
					"file": remotefile,
				}).Error("fetch: failed to read")
			} else {
				newBody := strings.Replace(string(fileData), " ", "", -1)
				var f interface{}
				err := json.Unmarshal([]byte(newBody), &f)
				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"name": cfg.Name,
						"file": remotefile,
					}).Error("fetch: failed to unmarshal")
				} else {
					*dataStore = append(*dataStore, f)
				}
			}

		}
	}

}

func getSSHConnection(yml *load.Config, api load.API) (*sftp.Client, error) {

	var user, pass, passphrase, host, port, sshpemfile string
	var encryptedpass, realpass []byte
	var err error
	var conn *ssh.Client
	var client *sftp.Client
	var authMethod ssh.AuthMethod
	var timeout time.Duration
	host = api.Scp.Host
	port = api.Scp.Port

	if yml.Global.SSHPEMFile != "" {
		sshpemfile = yml.Global.SSHPEMFile
	}
	if api.Scp.SSHPEMFile != "" {
		sshpemfile = api.Scp.SSHPEMFile
	}
	if yml.Global.User != "" {
		user = yml.Global.User
	}
	if yml.Global.Timeout > 0 {
		timeout = time.Duration(yml.Global.Timeout) * time.Millisecond
	} else {
		timeout = load.DefaultPingTimeout
	}

	if api.Scp.User != "" {
		user = api.Scp.User
	}

	if sshpemfile != "" {
		authMethod, err = publicKeyFile(sshpemfile)
	} else {

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
			encryptedpass, err = hex.DecodeString(pass)
			if err == nil {
				realpass, err = utils.Decrypt(encryptedpass, passphrase)
				if err == nil {
					pass = string(realpass)
				}
			}
		}
		err = nil
		authMethod = ssh.Password(pass)
	}

	if err == nil {

		sshconfig := &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         timeout,
			Auth: []ssh.AuthMethod{
				authMethod,
			},
		} // #nosec

		sshconfig.SetDefaults()
		conn, err = ssh.Dial("tcp", host+":"+port, sshconfig)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"user":  user,
				"host:": host,
				"err":   err,
			}).Error("ssh: failed to connect to host")
		} else {
			client, err = sftp.NewClient(conn)
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("ssh: failed to connect to sftp")
			}
			return client, err
		}
	}

	return nil, err
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"file": file,
			}).Error("Failed to read ssh pem file")
		}
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(key), err
}
