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
	"github.com/newrelic/nri-flex/internal/logger"
	yaml "gopkg.in/yaml.v2"
)

// loadSecrets if secrets configured fetch, store and substitute secrets
func loadSecrets(config *load.Config) {
	var ymlStr string

	for name, secret := range config.Secrets {
		if secret.Kind == "" {
			logger.Flex("error", fmt.Errorf("secret name: %v, missing kind", name), "", false)
			break
		}
		if secret.File == "" && secret.Data == "" && secret.HTTP.URL == "" {
			logger.Flex("error", fmt.Errorf("secret name: %v, file, data or http parameter needs to be set", name), "", false)
			break
		}

		logger.Flex("debug", nil, fmt.Sprintf("fetching secret for name: %v, kind: %v", name, secret.Kind), false)

		tempSecret := secret
		secretResult := ""

		switch secret.Kind {
		case "aws-kms":
			if secret.Region == "" {
				logger.Flex("error", fmt.Errorf("secret name: %v, missing region", secret.Region), "", false)
				break
			}

			secretResult = awskmsDecrypt(name, tempSecret)
		}

		if secretResult != "" {
			// convert config to string, only the first time
			if ymlStr == "" {
				ymlBytes, err := yaml.Marshal(config)
				if err != nil {
					logger.Flex("error", err, "", false)
					break
				}
				ymlStr = string(ymlBytes)
			}

			results := map[string]interface{}{
				"secret.result": secretResult,
			}

			if secret.Type != "" {
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
			logger.Flex("error", err, "", false)
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

// awskmsDecrypt perform aws kms decrypt and return plaintext
func awskmsDecrypt(name string, secret load.Secret) string {
	logger.Flex("debug", nil, "attempting to aws kms decrypt "+name+" secret", false)
	secretData := []byte{}

	if secret.File != "" {
		var fileData []byte
		fileData, err := ioutil.ReadFile(secret.File)
		if err == nil {
			secretData, err = base64.StdEncoding.DecodeString(string(fileData))
			if err != nil {
				logger.Flex("error", err, "", false)
			}
		} else {
			logger.Flex("error", err, "", false)
		}
	} else if secret.Data != "" {
		var err error
		secretData, err = base64.StdEncoding.DecodeString(secret.Data)
		logger.Flex("error", err, "", false)
	} else if secret.HTTP.URL != "" {
		client := &http.Client{}
		tlsConf := &tls.Config{}

		if secret.HTTP.TLSConfig.InsecureSkipVerify {
			tlsConf.InsecureSkipVerify = secret.HTTP.TLSConfig.InsecureSkipVerify
		}

		if secret.HTTP.TLSConfig.Ca != "" {
			rootCAs := x509.NewCertPool()
			ca, err := ioutil.ReadFile(secret.HTTP.TLSConfig.Ca)
			if err != nil {
				logger.Flex("error", err, "failed to read ca", false)
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
			logger.Flex("error", err, "", false)
		} else {
			for header, value := range secret.HTTP.Headers {
				req.Header.Add(header, value)
			}
			resp, err := client.Do(req)
			if err != nil {
				logger.Flex("error", err, "", false)
			} else {
				if resp.StatusCode == http.StatusOK {
					bodyBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						logger.Flex("error", err, "", false)
					} else {
						var err error
						secretData, err = base64.StdEncoding.DecodeString(string(bodyBytes))
						logger.Flex("error", err, "", false)
					}
				}
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
			logger.Flex("debug", nil, "aws kms decrypt "+name+" using custom credentials and/or config", false)
			sess = session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				SharedConfigFiles: sharedConfigFiles,
			}))
		} else {
			logger.Flex("debug", nil, "aws kms decrypt "+name+" using default credentials", false)
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
			logger.Flex("error", err, "aws kms decrypt "+name+" secret, fail", false)
		} else {
			logger.Flex("debug", nil, "aws kms decrypt "+name+" secret, success", false)
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
			logger.Flex("error", err, "", false)
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
