/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestCanRunMultipleCommands(t *testing.T) {
	// given
	load.Refresh()
	dataStoreExpected := []interface{}{
		map[string]interface{}{
			"active_defrag_hits":             "0",
			"active_defrag_key_hits":         "0",
			"active_defrag_key_misses":       "0",
			"active_defrag_misses":           "0",
			"active_defrag_running":          "0",
			"aof_current_rewrite_time_sec":   "-1",
			"aof_enabled":                    "0",
			"aof_last_bgrewrite_status":      "ok",
			"aof_last_cow_size":              "0",
			"aof_last_rewrite_time_sec":      "-1",
			"aof_last_write_status":          "ok",
			"aof_rewrite_in_progress":        "0",
			"aof_rewrite_scheduled":          "0",
			"arch_bits":                      "64",
			"atomicvar_api":                  "atomic-builtin",
			"blocked_clients":                "0",
			"client_biggest_input_buf":       "0",
			"client_longest_output_list":     "0",
			"cluster_enabled":                "0",
			"config_file":                    "/usr/local/etc/redis.conf",
			"connected_clients":              "1",
			"connected_slaves":               "0",
			"db0":                            "keys=1,expires=0,avg_ttl=0",
			"evicted_keys":                   "0",
			"executable":                     "/usr/local/opt/redis/bin/redis-server",
			"expired_keys":                   "0",
			"expired_stale_perc":             "0.00",
			"expired_time_cap_reached_count": "0",
			"gcc_version":                    "4.2.1",
			"hz":                             "10",
			"instantaneous_input_kbps":       "0.00",
			"instantaneous_ops_per_sec":      "0",
			"instantaneous_output_kbps":      "0.00",
			"keyspace_hits":                  "1",
			"keyspace_misses":                "1",
			"latest_fork_usec":               "1234",
			"lazyfree_pending_objects":       "0",
			"loading":                        "0",
			"lru_clock":                      "6777537",
			"master_repl_offset":             "0",
			"master_replid":                  "65a558faf348fb6f0cf75b82e4ae6b1cc4128baf",
			"master_replid2":                 "0000000000000000000000000000000000000000",
			"maxmemory":                      "0",
			"maxmemory_human":                "0B",
			"maxmemory_policy":               "noeviction",
			"mem_allocator":                  "libc",
			"mem_fragmentation_ratio":        "0.53",
			"migrate_cached_sockets":         "0",
			"multiplexing_api":               "kqueue",
			"os":                             "Darwin 18.2.0 x86_64",
			"process_id":                     "1201",
			"pubsub_channels":                "0",
			"pubsub_patterns":                "0",
			"rdb_bgsave_in_progress":         "0",
			"rdb_changes_since_last_save":    "0",
			"rdb_current_bgsave_time_sec":    "-1",
			"rdb_last_bgsave_status":         "ok",
			"rdb_last_bgsave_time_sec":       "0",
			"rdb_last_cow_size":              "0",
			"rdb_last_save_time":             "1548333171",
			"redis_build_id":                 "e0c8d37381c486c6",
			"redis_git_dirty":                "0",
			"redis_git_sha1":                 "00000000",
			"redis_mode":                     "standalone",
			"redis_version":                  "4.0.9",
			"rejected_connections":           "0",
			"repl_backlog_active":            "0",
			"repl_backlog_first_byte_offset": "0",
			"repl_backlog_histlen":           "0",
			"repl_backlog_size":              "1048576",
			"role":                           "master",
			"run_id":                         "7d1e18e7ccf9484ef3636ff145dbc3ced5fd8f01",
			"second_repl_offset":             "-1",
			"slave_expires_tracked_keys":     "0",
			"sync_full":                      "0",
			"sync_partial_err":               "0",
			"sync_partial_ok":                "0",
			"tcp_port":                       "6379",
			"total_commands_processed":       "331",
			"total_connections_received":     "323",
			"total_net_input_bytes":          "0",
			"total_net_output_bytes":         "871587",
			"total_system_memory":            "17179869184",
			"total_system_memory_human":      "16.00G",
			"uptime_in_days":                 "38",
			"uptime_in_seconds":              "3349376",
			"used_cpu_sys":                   "422.27",
			"used_cpu_sys_children":          "0.00",
			"used_cpu_user":                  "200.05",
			"used_cpu_user_children":         "0.00",
			"used_memory":                    "1043984",
			"used_memory_dataset":            "13546",
			"used_memory_dataset_perc":       "21.42%",
			"used_memory_human":              "1019.52K",
			"used_memory_lua":                "37888",
			"used_memory_lua_human":          "37.00K",
			"used_memory_overhead":           "1030438",
			"used_memory_peak":               "1043984",
			"used_memory_peak_human":         "1019.52K",
			"used_memory_peak_perc":          "100.00%",
			"used_memory_rss":                "557056",
			"used_memory_rss_human":          "544.00K",
			"used_memory_startup":            "980736",
			"zHost":                          "HELLO",
		},
	}

	configFile := load.Config{
		Name: "RedisInfo",
		APIs: []load.API{
			{
				Name:     "redis",
				Commands: getCanRunMultipleCommands(),
			},
		},
	}

	// when
	dataStore := []interface{}{}
	RunCommands(&dataStore, &configFile, 0)

	actual := dataStore[0].(map[string]interface{})
	expected := dataStoreExpected[0].(map[string]interface{})
	// then each of the values found in the expected is equal to the actual
	for key, expectedValue := range expected {
		if key == "flex.commandTimeMs" {
			continue
		}
		actualValue := actual[key]
		assert.Equalf(t, expectedValue, actualValue, "%s doesnt match - want: %v  got: %v", key, expectedValue, actualValue)
	}
}

func TestDf(t *testing.T) {
	load.Refresh()
	config := load.Config{
		Name: "dfFlex",
		APIs: getDfApis(),
	}

	dataStoreExpected := []interface{}{
		map[string]interface{}{
			"mountedOn": "/",
			"used":      "744562808",
			"available": "224223192",
			"capacity":  "77%",
			"ifree":     "9223372036851372364",
			"iusedPerc": "0%",
			"fs":        "/dev/disk1s1",
			"512Blocks": "976490568",
			"iused":     "3403443",
		},
	}

	dataStore := []interface{}{}
	RunCommands(&dataStore, &config, 0)

	assert.Len(t, dataStore, 3)

	// we are only checking the first entry
	expected := dataStoreExpected[0].(map[string]interface{})
	actual := dataStore[0].(map[string]interface{})

	for key, expectedValue := range expected {
		actualValue := actual[key]
		assert.Equalf(t, expectedValue, actualValue, "%s doesnt match - want: %v  got: %v", key, expectedValue, actualValue)
	}
}

func TestDf2(t *testing.T) {
	load.Refresh()
	config := load.Config{
		Name: "dfFlex",
		APIs: getDf2Apis(),
	}

	dataStoreExpected := []interface{}{
		map[string]interface{}{
			"Mounted":    "/",
			"Used":       "744562808",
			"Available":  "224223192",
			"Capacity":   "77",
			"ifree":      "9223372036851372364",
			"%iused":     "0",
			"Filesystem": "/dev/disk1s1",
			"512-blocks": "976490568",
			"iused":      "3403443",
		},
	}

	dataStore := []interface{}{}
	RunCommands(&dataStore, &config, 0)

	assert.Len(t, dataStore, 3)

	// we are only checking the first entry
	expected := dataStoreExpected[0].(map[string]interface{})
	actual := dataStore[0].(map[string]interface{})

	for key, expectedValue := range expected {
		actualValue := actual[key]
		assert.Equal(t, expectedValue, actualValue)
	}
}

func TestRawCache(t *testing.T) {
	load.Refresh()
	config := load.Config{
		RawCache: map[string]interface{}{},
		Name:     "rawCacheExample",
		APIs:     getRawCacheApis(),
	}

	dataStoreExpected := []interface{}{
		map[string]interface{}{
			"batman": "bruce",
		},
	}

	dataStore := []interface{}{}
	RunCommands(&dataStore, &config, 0)
	RunCommands(&dataStore, &config, 1)

	assert.Len(t, dataStore, 1)

	// we are only checking the first entry
	expected := dataStoreExpected[0].(map[string]interface{})
	actual := dataStore[0].(map[string]interface{})

	for key, expectedValue := range expected {
		actualValue := actual[key]
		assert.Equal(t, expectedValue, actualValue)
	}
}
