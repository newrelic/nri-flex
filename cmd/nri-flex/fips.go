// Copyright 2019 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build fips
// +build fips

package main

import (
	_ "crypto/tls/fipsonly"
)
