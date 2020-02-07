/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
)

// SetJMXCommand Add parameters to JMX call
func SetJMXCommand(dataStore *[]interface{}, runCommand *string, command load.Command, api load.API, config *load.Config) {
	*runCommand = fmt.Sprintf("echo '%v' | java -jar %vnrjmx.jar", *runCommand, load.Args.NRJMXToolPath)

	// order command > api > global
	if command.Jmx.Host != "" {
		*runCommand = *runCommand + " -hostname " + command.Jmx.Host
	} else if api.Jmx.Host != "" {
		*runCommand = *runCommand + " -hostname " + api.Jmx.Host
	} else if config.Global.Jmx.Host != "" {
		*runCommand = *runCommand + " -hostname " + config.Global.Jmx.Host
	} else {
		*runCommand = *runCommand + " -hostname " + load.DefaultJmxHost
	}

	if command.Jmx.Port != "" {
		*runCommand = *runCommand + " -port " + command.Jmx.Port
	} else if api.Jmx.Port != "" {
		*runCommand = *runCommand + " -port " + api.Jmx.Port
	} else if config.Global.Jmx.Port != "" {
		*runCommand = *runCommand + " -port " + config.Global.Jmx.Port
	}

	if command.Jmx.User != "" {
		*runCommand = *runCommand + " -username " + command.Jmx.User
	} else if api.Jmx.User != "" {
		*runCommand = *runCommand + " -username " + api.Jmx.User
	} else if config.Global.Jmx.User != "" {
		*runCommand = *runCommand + " -username " + config.Global.Jmx.User
	}

	if command.Jmx.Pass != "" {
		*runCommand = *runCommand + " -password " + command.Jmx.Pass
	} else if api.Jmx.Pass != "" {
		*runCommand = *runCommand + " -password " + api.Jmx.Pass
	} else if config.Global.Jmx.Pass != "" {
		*runCommand = *runCommand + " -password " + config.Global.Jmx.Pass
	}

	if command.Jmx.KeyStore != "" {
		*runCommand = *runCommand + " -keyStore " + command.Jmx.KeyStore
	} else if api.Jmx.KeyStore != "" {
		*runCommand = *runCommand + " -keyStore " + api.Jmx.KeyStore
	} else if config.Global.Jmx.KeyStore != "" {
		*runCommand = *runCommand + " -keyStore " + config.Global.Jmx.KeyStore
	}

	if command.Jmx.KeyStorePass != "" {
		*runCommand = *runCommand + " -keyStorePassword " + command.Jmx.KeyStorePass
	} else if api.Jmx.KeyStorePass != "" {
		*runCommand = *runCommand + " -keyStorePassword " + api.Jmx.KeyStorePass
	} else if config.Global.Jmx.KeyStorePass != "" {
		*runCommand = *runCommand + " -keyStorePassword " + config.Global.Jmx.KeyStorePass
	}

	if command.Jmx.TrustStore != "" {
		*runCommand = *runCommand + " -trustStore " + command.Jmx.TrustStore
	} else if api.Jmx.TrustStore != "" {
		*runCommand = *runCommand + " -trustStore " + api.Jmx.TrustStore
	} else if config.Global.Jmx.TrustStore != "" {
		*runCommand = *runCommand + " -trustStore " + config.Global.Jmx.TrustStore
	}

	if command.Jmx.TrustStorePass != "" {
		*runCommand = *runCommand + " -trustStorePassword " + command.Jmx.TrustStorePass
	} else if api.Jmx.TrustStorePass != "" {
		*runCommand = *runCommand + " -trustStorePassword " + api.Jmx.TrustStorePass
	} else if config.Global.Jmx.TrustStorePass != "" {
		*runCommand = *runCommand + " -trustStorePassword " + config.Global.Jmx.TrustStorePass
	}

	if command.Jmx.URIPath != "" {
		*runCommand = *runCommand + " -uriPath " + command.Jmx.URIPath
	} else if api.Jmx.URIPath != "" {
		*runCommand = *runCommand + " -uriPath " + api.Jmx.URIPath
	} else if config.Global.Jmx.URIPath != "" {
		*runCommand = *runCommand + " -uriPath " + config.Global.Jmx.URIPath
	}

	load.Logrus.Debugf("commands: completed jmx command: %v", *runCommand)
}

// ParseJMX Processes JMX Data
func ParseJMX(dataStore *[]interface{}, dataInterface interface{}, command load.Command, dataSample *map[string]interface{}) {
	// dataSample contains data from previously run raw commands

	sendSample := true
	data := dataInterface.(map[string]interface{})
	if command.CompressBean {
		newJMXSample := map[string]interface{}{}
		for key, val := range data {
			compressedBean := compressBean(key)
			newJMXSample[compressedBean] = val
			delete(data, key)
		}
		*dataStore = append(*dataStore, newJMXSample)
		// load.StoreAppend(newJMXSample)
	} else {
		jmxSamples := map[string]map[string]interface{}{}
		for k, v := range data {
			keyArray := strings.Split(k, ",")
			groupKey := ""

			if command.GroupBy == "" {
				groupKey = getBeanName(k)
			} else {
				// find group key first
				for _, k2 := range keyArray {
					keyArray2 := strings.Split(k2, "=")
					if len(keyArray2) == 2 {
						if keyArray2[0] == command.GroupBy {
							groupKey = keyArray2[1]
							break
						}
					}
				}
			}

			if jmxSamples[groupKey] == nil {
				jmxSamples[groupKey] = map[string]interface{}{}
				domain, query := splitBeanName(groupKey)
				jmxSamples[groupKey]["bean"] = query
				jmxSamples[groupKey]["domain"] = domain

				// add raw command data from dataSample
				for k, v := range *dataSample {
					jmxSamples[groupKey][k] = v
				}
			}

			for _, k2 := range keyArray {
				keyArray2 := strings.Split(k2, "=")
				if len(keyArray2) == 2 {
					if keyArray2[0] == "attr" {
						jmxSamples[groupKey][keyArray2[1]] = v
					} else {
						// use key filters here
						// ?? dual regex map[string]string
						jmxSamples[groupKey][keyArray2[0]] = keyArray2[1]
					}
				}
			}
		}

		for _, jmxSample := range jmxSamples {
			if sendSample {
				*dataStore = append(*dataStore, jmxSample)
				// load.StoreAppend(jmxSample)
			}
		}
	}

}

func getBeanName(beanString string) string {
	beanNameRegex := regexp.MustCompile("^(.*),attr=.*")
	beanNameMatches := beanNameRegex.FindStringSubmatch(beanString)
	if beanNameMatches == nil {
		return "FailedToGetBeanName"
	}
	return beanNameMatches[1]
}

func getAttrName(beanString string) string {
	attrNameRegex := regexp.MustCompile("^.*attr=(.*)$")
	attrNameMatches := attrNameRegex.FindStringSubmatch(beanString)
	if attrNameMatches == nil {
		// fmt.Println(beanString, attrNameRegex)
		return "FailedToGetAttrName"
	}
	return attrNameMatches[1]
}

func splitBeanName(bean string) (string, string) {
	domainQuery := strings.SplitN(bean, ":", 2)
	if len(domainQuery) != 2 {
		return "CouldNotGetDomain", "CouldNotGetQuery"
	}
	return domainQuery[0], domainQuery[1]
}

func compressBean(bean string) string {
	attr := getAttrName(bean)
	fullBean := getBeanName(bean)
	splitBean := strings.SplitN(fullBean, ":", 2)
	newBean := ""
	if len(splitBean) == 2 {
		splitBeans := strings.Split(splitBean[1], ",")
		for _, str := range splitBeans {
			key := strings.Split(str, "=")
			if len(key) == 2 {
				if newBean == "" {
					newBean = key[1]
				} else {
					newBean = newBean + "." + key[1]
				}
			}
		}
	}
	// if attr == "FailedToGetAttrName" {
	// 	fmt.Println(attr, " ::: ", bean, " ::: ", newBean)
	// } else {
	// 	fmt.Println("GOT IT             ", " ::: ", bean, " ::: ", newBean)
	// }
	return newBean + "." + attr
}

// func getKeyProperties(keyProperties string) (map[string]string, error) {
// 	keyPropertiesMap := map[string]string{}
// 	keyPropertiesArray := strings.Split(keyProperties, ",")
// 	for _, keyProperty := range keyPropertiesArray {
// 		keyPropertySplit := strings.Split(keyProperty, "=")
// 		if len(keyPropertySplit) != 2 {
// 			return nil, fmt.Errorf("invalid key properties %s", keyProperties)
// 		}
// 		keyPropertiesMap[keyPropertySplit[0]] = keyPropertySplit[1]
// 	}
// 	return keyPropertiesMap, nil
// }
