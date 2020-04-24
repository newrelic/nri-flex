//+build windows

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
			Run:     "type ..\\..\\test\\payloads\\redisInfo.out",
			SplitBy: ":",
		},
		{
			Run:     "echo zHost:HELLO",
			SplitBy: ":",
		},
	}
}

func getDfApis() []load.API {
	return []load.API{
		{
			Name: "df",
			Commands: []load.Command{
				{
					Run:      "type ..\\..\\test\\payloads\\df.out",
					Split:    "horizontal",
					RowStart: 1,
					SetHeader: []string{
						"fs", "512Blocks", "used", "available", "capacity", "iused", "ifree", "iusedPerc", "mountedOn",
					},
					RegexMatch: false,
					SplitBy:    `\s{1,}`,
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
					Run:              "type ..\\..\\test\\payloads\\df.out",
					Split:            "horizontal",
					RegexMatch:       true,
					SplitBy:          `(\S+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(\d+)\s+(\d+)\s+(\d+)%\s+(\W*)`,
					HeaderRegexMatch: false,
					HeaderSplitBy:    `\s{1,}`,
				},
			},
		},
	}
}
