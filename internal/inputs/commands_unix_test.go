// +build linux darwin

/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"github.com/newrelic/nri-flex/internal/load"
)

func getCanRunMultipleCommands() []load.Command {
	return []load.Command{
		{
			Run:     "cat ../../test/payloads/redisInfo.out",
			SplitBy: ":",
		},
		{
			Run:     `echo "zHost:$(echo HELLO)"`,
			SplitBy: ":",
		},
	}
}

func getDfApis() []load.API {
	return []load.API{
		{
			Name: "df",
			//Shell: "/bin/sh",
			Commands: []load.Command{
				{
					Run:      "cat ../../test/payloads/df.out",
					Split:    "horizontal",
					RowStart: 1,
					SetHeader: []string{
						"fs", "512Blocks", "used", "available", "capacity", "iused", "ifree", "iusedPerc", "mountedOn",
					},
					RegexMatch: false,
					SplitBy:    `\s{1,}`,
					//Shell:      "/bin/sh",
				},
			},
		},
	}
}

func getDf2Apis() []load.API {
	return []load.API{
		{
			Name: "df",
			Commands: []load.Command{
				{
					Run:              "cat ../../test/payloads/df.out",
					Split:            "horizontal",
					RegexMatch:       true,
					SplitBy:          `(\S+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(\d+)\s+(\d+)\s+(\d+)%\s+(.*)`,
					HeaderRegexMatch: false,
					HeaderSplitBy:    `\s{1,}`,
				},
			},
		},
	}
}

func getRawCacheApis() []load.API {
	return []load.API{
		{
			Name: "getSomeData",
			Commands: []load.Command{
				{
					Name:         "info",
					Run:          "echo batman:bruce",
					IgnoreOutput: true,
				},
			},
		},
		{
			Name: "hero",
			Commands: []load.Command{
				{
					Run:     "echo ${cache:info}",
					SplitBy: ":",
				},
			},
		},
	}
}
