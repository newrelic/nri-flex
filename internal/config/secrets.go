/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/utils"
)

// loadSecrets if secrets configured fetch, store and substitute secrets
func loadSecrets(config *load.Config) error {
	var ymlStr string
	var err error
	for name, secret := range config.Secrets {
		if secret.Kind == "" {
			err = fmt.Errorf("config: secret needs 'kind' parameter to be set")
			load.Logrus.WithFields(logrus.Fields{
				"secret": name,
			}).Error(err.Error())
			break
		}
		if secret.File == "" && secret.Data == "" && secret.HTTP.URL == "" {
			err = fmt.Errorf("config: secret needs 'file', 'data' and 'http' parameter to be set")
			load.Logrus.WithFields(logrus.Fields{
				"secret": name,
			}).Error(err.Error())
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
				err = fmt.Errorf("config: secret needs 'region' parameter to be set")
				load.Logrus.WithFields(logrus.Fields{
					"secret": name,
					"kind":   secret.Kind,
				}).Error(err.Error())
				break
			}
			secretResult = awskmsDecrypt(name, tempSecret)
		case "vault":
			if secret.HTTP.URL == "" {
				err = fmt.Errorf("config: vault secret requires 'http' parameter to be set")
				load.Logrus.WithFields(logrus.Fields{
					"secret": name,
					"kind":   secret.Kind,
				}).Error(err.Error())
				break
			}
			vaultFetch(name, tempSecret, results)
			// decrypt secret locally using simpleEncrypDecryp module
		case "local":
			if secret.Key == "" {
				err = fmt.Errorf("config: local secret requires 'key' parameter to be set")
				load.Logrus.WithFields(logrus.Fields{
					"secret": name,
					"kind":   secret.Kind,
				}).Error(err.Error())
				break
			}
			secretResult = localDecrypt(name, tempSecret)
		}

		if secretResult != "" || len(results) > 0 {
			// convert config to string, only the first time
			if ymlStr == "" {
				ymlBytes, e := yaml.Marshal(config)
				if e != nil {
					err = fmt.Errorf("config: secret marshal failed")
					load.Logrus.WithFields(logrus.Fields{
						"secret": name,
						"kind":   secret.Kind,
						"err":    e,
					}).Error(err.Error())
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
		*config, err = ReadYML(ymlStr)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"config": config.Name,
				"err":    err,
			}).Error("config: secret unmarshal failed")
		}
	}
	return err
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
	var secretData []byte

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
			return ""
		}
		result := string(resp.Plaintext)
		load.Logrus.WithFields(logrus.Fields{
			"name":   name,
			"secret": result,
		}).Debug("config: aws kms decrypt secret success")
		return result
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
			load.Logrus.WithError(err).Error("config: secret failed to read tls ca")
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

// localDecrypt perform local decrypt and return plaintext if decrpyted successfully
func localDecrypt(name string, secret load.Secret) string {
	load.Logrus.WithFields(logrus.Fields{"name": name}).Debug("config: attempting to local decrypt secret")
	var secretData []byte

	if secret.File != "" {
		var fileData []byte
		fileData, err := ioutil.ReadFile(secret.File)
		if err == nil {
			secretData, err = hex.DecodeString(string(fileData))
			if err != nil {
				load.Logrus.WithFields(logrus.Fields{
					"name": name,
				}).WithError(err).Error("config: local secret hex decode failed")
			}
		} else {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
				"file": secret.File,
			}).WithError(err).Error("config: local read file failed")
		}
	} else if secret.Data != "" {
		var err error
		secretData, err = hex.DecodeString(secret.Data)
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
			}).WithError(err).Error("config: local secret hex decode failed")
		}
	}

	if len(secretData) > 0 {
		if secret.Key != "" {
			result, err := utils.Decrypt(secretData, secret.Key)
			if err == nil {
				return string(result)
			}
			load.Logrus.WithFields(logrus.Fields{
				"name": name,
				"key":  secret.Key,
			}).Error("config: local unable to decrypt using key provided, return encrypted data as is")

		}

	}
	return string(secretData)
}
