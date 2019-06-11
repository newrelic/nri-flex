package processor

import "fmt"

// StoreLookups if key is found (using regex), store the values in the lookupStore as the defined lookupStoreKey for later use
func StoreLookups(storeLookups map[string]string, key *string, lookupStore *map[string][]string, v *interface{}) {
	for lookupStoreKey, lookupFindKey := range storeLookups {
		if *key == lookupFindKey {
			if (*lookupStore) == nil {
				(*lookupStore) = map[string][]string{}
			}
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
