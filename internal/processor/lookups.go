/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package processor

import (
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

func cleanValue(v *interface{}) string {
	switch val := (*v).(type) {
	case float32, float64:
		return fmt.Sprintf("%f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// StoreLookups if key is found (using regex), store the values in the lookupStore as the defined lookupStoreKey for later use
func StoreLookups(storeLookups map[string]string, key *string, lookupStore *map[string]map[string]struct{}, v *interface{}) {
	for lookupStoreKey, lookupFindKey := range storeLookups {
		if *key == lookupFindKey {
			load.Logrus.WithFields(logrus.Fields{
				"lookupFindKey": lookupFindKey,
				lookupStoreKey:  fmt.Sprintf("%v", *v),
			}).Debug("create: store lookup")

			if (*lookupStore)[lookupStoreKey] == nil {
				(*lookupStore)[lookupStoreKey] = make(map[string]struct{})
			}

			switch data := (*v).(type) {
			case []interface{}:
				load.Logrus.WithFields(logrus.Fields{
					"lookupFindKey": lookupFindKey,
				}).Debug("splitting array")

				for _, dataKey := range data {
					(*lookupStore)[lookupStoreKey][cleanValue(&dataKey)] = struct{}{}
				}
			default:
				(*lookupStore)[lookupStoreKey][cleanValue(v)] = struct{}{}
			}
		}
	}
}

// VariableLookups if key is found (using regex), store the value in the variableStore, as the defined by the variableStoreKey for later use
func VariableLookups(variableLookups map[string]string, key *string, variableStore *map[string]string, v *interface{}) {
	for variableStoreKey, variableFindKey := range variableLookups {
		if *key == variableFindKey {
			if (*variableStore) == nil {
				(*variableStore) = map[string]string{}
			}
			(*variableStore)[variableStoreKey] = cleanValue(v)
		}
	}
}
