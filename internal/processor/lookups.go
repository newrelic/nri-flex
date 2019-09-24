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

// StoreLookups if key is found (using regex), store the values in the lookupStore as the defined lookupStoreKey for later use
func StoreLookups(storeLookups map[string]string, key *string, lookupStore *map[string][]string, v *interface{}) {
	for lookupStoreKey, lookupFindKey := range storeLookups {
		if *key == lookupFindKey {
			load.Logrus.WithFields(logrus.Fields{
				"lookupFindKey": lookupFindKey,
				lookupStoreKey:  fmt.Sprintf("%v", *v),
			}).Debug("create: store lookup")

			(*lookupStore)[lookupStoreKey] = append((*lookupStore)[lookupStoreKey], fmt.Sprintf("%v", *v))
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
			(*variableStore)[variableStoreKey] = fmt.Sprintf("%v", *v)
		}
	}
}
