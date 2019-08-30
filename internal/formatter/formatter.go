/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package formatter

import (
	"fmt"
	"regexp"
	"strings"
)

// ValueParse Plucks first found value out with regex, if nothing found send back the value
func ValueParse(v interface{}, regex string) string {
	value := fmt.Sprintf("%v", v)
	reg := regexp.MustCompile(regex)
	matches := reg.FindAllString(value, -1)
	if len(matches) > 0 {
		return matches[0] // send first match back
	}
	return value
}

// SplitKey simple key value pair splitter
func SplitKey(key, splitChar string) (string, string, bool) {
	keys := strings.SplitN(key, splitChar, 2)
	if len(keys) == 2 {
		return keys[0], keys[1], true
	}
	return "", "", false
}

// PercToDecimal convert percentage to decimal
func PercToDecimal(v *interface{}) {
	*v = strings.TrimRight(fmt.Sprintf("%v", *v), "%")
}

// SnakeCaseToCamelCase converts snake_case to camelCase
func SnakeCaseToCamelCase(key *string) {
	isToUpper := false
	for k, v := range *key {
		if k == 0 {
			*key = strings.ToLower(string((*key)[0]))
		} else {
			if isToUpper {
				*key += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					*key += string(v)
				}
			}
		}
	}
}

// RegMatch Perform regex matching
func RegMatch(text string, regexmatch string) []string {
	reg := regexp.MustCompile(regexmatch)
	matches := reg.FindStringSubmatch(text)
	if matches != nil {
		return matches[1:]
	}
	return nil
}

// RegSplit Split by Regex
func RegSplit(text string, delimiter string) []string {
	reg := regexp.MustCompile(delimiter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:]
	return result
}

// KvFinder tests with multiple modes, whether k1 satisfies k2
func KvFinder(mode string, k1 string, k2 string) bool {
	switch {
	case mode == "prefix" && strings.HasPrefix(k1, k2):
		return true
	case mode == "suffix" && strings.HasSuffix(k1, k2):
		return true
	case mode == "contains" && strings.Contains(k1, k2):
		return true
	case mode == "regex":
		validateKey := regexp.MustCompile(k2)
		return validateKey.MatchString(k1)
	default:
		return false
	}
}
