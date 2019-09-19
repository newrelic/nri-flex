/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// loadSecrets if secrets configured fetch, store and substitute secrets
func loadSecrets(config *load.Config) {
	var ymlStr string

	for name, secret := range config.Secrets {
		if secret.Kind == "" {
			load.Logrus.WithFields(logrus.Fields{
				"secret": name,
			}).Error("config: secret is missing kind")
			break
		}
		if secret.File == "" && secret.Data == "" && secret.HTTP.URL == "" {
			load.Logrus.WithFields(logrus.Fields{
				"secret": name,
			}).Error(fmt.Sprintf("config: secret needs file, data or http parameter needs to be set"))
			break
		}

		load.Logrus.WithFields(logrus.Fields{
			"secret": name,
			"kind":   secret.Kind,
		}).Debug("config: fetching secret")

		tempSecret := secret
		secretResult := ""
		results := map[string]interface{}{}

		switch secret.Kind {
		case "aws-kms":
			if secret.Region == "" {
				load.Logrus.WithFields(logrus.Fields{
					"secret": name,
					"kind":   secret.Kind,
				}).Error("config: secret missing region")
				break
			}
			secretResult = awskmsDecrypt(name, tempSecret)
		case "vault":
			if secret.HTTP.URL == "" {
				load.Logrus.WithFields(logrus.Fields{
					"secret": name,
					"kind":   secret.Kind,
				}).Error("config: vault secret requires http parameter to be set")
				break
			}
			vaultFetch(name, tempSecret, results)
		}

		if secretResult != "" || len(results) > 0 {
			// convert config to string, only the first time
			if ymlStr == "" {
				ymlBytes, err := yaml.Marshal(config)
				if err != nil {
					load.Logrus.WithFields(logrus.Fields{
						"secret": name,
						"kind":   secret.Kind,
						"err":    err,
					}).Error("config: secret marshal failed")
					break
				}
				ymlStr = string(ymlBytes)
			}

			if secretResult != "" {
				results["secret.result"] = secretResult
			}

			if secret.Type != "" && secret.Kind != "vault" {
				handleDataType(results, secret.Type)
			}

			ymlStr = subSecrets(ymlStr, name, results)
		}
	}

	// if ymlStr has a value it means a secret was successfully retrieved, decrypted, and substitutions were attempted
	// we can then attempt to read and overwrite the config
	if ymlStr != "" {
		var err error
		*config, err = ReadYML(ymlStr)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"config": config.Name,
				"err":    err,
			}).Error("config: secret unmarshal failed")
		}
	}
}

// subSecrets substitute secrets into yml str and return
func subSecrets(configStr string, secretKey string, secrets map[string]interface{}) string {
	variableReplaces := regexp.MustCompile(`\${secret\.`+secretKey+`:.*?}`).FindAllString(configStr, -1)
	for _, variableReplace := range variableReplaces {
		variableKey := strings.TrimSuffix(strings.Split(variableReplace, "${secret."+secretKey+":")[1], "}") // eg. "channel"
		if variableKey == "value" {
			configStr = strings.Replace(configStr, variableReplace, fmt.Sprintf("%v", secrets["secret.result"]), -1)
		} else if secrets[variableKey] != nil {
			configStr = strings.Replace(configStr, variableReplace, fmt.Sprintf("%v", secrets[variableKey]), -1)
		}
	}
	return configStr
}

// vaultFetch fetch from Hashicorp Vault
func vaultFetch(name string, secret load.Secret, results map[string]interface{}) {
	load.Logrus.WithFields(logrus.Fields{"name": name}).Debug("config: fetching vault secret")
	bytes, err := httpWrapper(secret)
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{"name": name, "err": err}).Error("config: fetching vault secret failed")
	} else {
		var jsonInterface map[string]interface{}
		err := json.Unmarshal(bytes, &jsonInterface)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{"name": name, "err": err}).Error("config: vault data unmarshal failed")
		} else {
			// v1 and v2 engines have this available
			if jsonInterface["data"] != nil {
				load.Logrus.WithFields(logrus.Fields{"name": name}).Debug("config: fetching vault secret success")
				switch firstData := jsonInterface["data"].(type) {
				case map[string]interface{}:
					isV2 := false
					if firstData["data"] != nil { // v2 format
						switch secondData := firstData["data"].(type) {
						case map[string]interface{}:
							// handle v2 data
							isV2 = true
							for key, val := range secondData {
								results[key] = val
							}
						}
					}
					if !isV2 {
						for key, val := range firstData {
							results[key] = val
						}
					}
				}
			}
		}
	}
}

// awskmsDecrypt perform aws kms decrypt and return plaintext
func awskmsDecrypt(name string, secret load.Secret) string {
	load.Logrus.WithFields(logrus.Fields{"name": name}).Debug("config: attempting to aws kms decrypt secret")
	secretData := []byte{}

	if secret.File != "" {
		var fileData []byte
		fileData, err := ioutil.ReadFile(secret.File)
		if err == nil {
			secretData, err = base64.StdEncoding.DecodeString(string(fileData))
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": name,
					"err":  err,
				}).Error("config: secret base64 decode failed")
			}
		} else {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
				"err":  err,
				"file": secret.File,
			}).Error("config: aws kms read file failed")
		}
	} else if secret.Data != "" {
		var err error
		secretData, err = base64.StdEncoding.DecodeString(secret.Data)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
				"err":  err,
			}).Error("config: aws kms base64 decode failed")
		}
	} else if secret.HTTP.URL != "" {
		bytes, err := httpWrapper(secret)
		if err != nil {

			load.Logrus.WithFields(logrus.Fields{
				"url":  secret.HTTP.URL,
				"name": name,
				"err":  err,
			}).Error("config: aws kms http fetch failed")

		} else {
			var err error
			secretData, err = base64.StdEncoding.DecodeString(string(bytes))
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": name,
					"err":  err,
				}).Error("config: aws kms base64 decode failed")
			}
		}
	}

	if len(secretData) > 0 {
		var sess *session.Session

		sharedConfigFiles := []string{}
		if secret.CredentialFile != "" {
			sharedConfigFiles = append(sharedConfigFiles, secret.CredentialFile)
		}
		if secret.ConfigFile != "" {
			sharedConfigFiles = append(sharedConfigFiles, secret.ConfigFile)
		}

		if len(sharedConfigFiles) > 0 {

			load.Logrus.WithFields(logrus.Fields{
				"name": name,
			}).Debug("config: aws kms decrypt using custom credentials and/or config")

			sess = session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				SharedConfigFiles: sharedConfigFiles,
			}))
		} else {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
			}).Debug("config: aws kms decrypt using default credentials")
			sess = session.Must(session.NewSession(&aws.Config{
				Region: aws.String(secret.Region),
			}))
		}

		kmsClient := kms.New(sess)
		params := &kms.DecryptInput{
			CiphertextBlob: secretData,
		}
		resp, err := kmsClient.Decrypt(params)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
			}).Error("config: aws kms decrypt secret failed")
		} else {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
			}).Error("config: aws kms decrypt secret success")
			return string(resp.Plaintext)
		}
	}
	return ""
}

func handleDataType(results map[string]interface{}, dataType string) {
	switch dataType {
	case "json":
		var jsonResult map[string]interface{}
		err := json.Unmarshal([]byte(results["secret.result"].(string)), &jsonResult)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("config: secret unmarshal failed")
		} else {
			for key, value := range jsonResult {
				results[key] = fmt.Sprintf("%v", value)
			}
		}
	case "equal":
		commaSplit := strings.Split(results["secret.result"].(string), ",")
		for _, initialSplit := range commaSplit {
			equalSplit := strings.SplitN(initialSplit, "=", 2)
			if len(equalSplit) == 2 {
				results[equalSplit[0]] = fmt.Sprintf("%v", equalSplit[1])
			}
		}
	}
}

func httpWrapper(secret load.Secret) ([]byte, error) {
	client := &http.Client{}
	tlsConf := &tls.Config{}

	if secret.HTTP.TLSConfig.InsecureSkipVerify {
		tlsConf.InsecureSkipVerify = secret.HTTP.TLSConfig.InsecureSkipVerify
	}

	if secret.HTTP.TLSConfig.Ca != "" {
		rootCAs := x509.NewCertPool()
		ca, err := ioutil.ReadFile(secret.HTTP.TLSConfig.Ca)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("config: secret failed to read ca")
		} else {
			rootCAs.AppendCertsFromPEM(ca)
			tlsConf.RootCAs = rootCAs
		}
	}

	clientConf := &http.Transport{
		TLSClientConfig: tlsConf,
	}

	client.Transport = clientConf
	req, err := http.NewRequest("GET", secret.HTTP.URL, nil)

	if err != nil {
		return nil, err
	}

	for header, value := range secret.HTTP.Headers {
		req.Header.Add(header, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return bytes, nil
	}

	return nil, fmt.Errorf("http fetch failed %v %v", resp.StatusCode, string(bytes))
}
