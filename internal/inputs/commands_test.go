package inputs

import (
	"fmt"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
)

func TestRedis(t *testing.T) {
	dataStoreExpected := []interface{}{
		map[string]interface{}{
			"redis_git_dirty":            0,
			"multiplexing_api":           "kqueue",
			"maxmemory":                  0,
			"instantaneous_output_kbps":  "0.00",
			"master_replid2":             "0000000000000000000000000000000000000000",
			"used_cpu_user":              200.05,
			"client_biggest_input_buf":   0,
			"total_system_memory":        17179869184,
			"mem_fragmentation_ratio":    0.53,
			"total_commands_processed":   331,
			"keyspace_hits":              1,
			"process_id":                 1201,
			"used_memory":                1043984,
			"instantaneous_input_kbps":   "0.00",
			"active_defrag_key_misses":   0,
			"master_replid":              "65a558faf348fb6f0cf75b82e4ae6b1cc4128baf",
			"lazyfree_pending_objects":   0,
			"rdb_bgsave_in_progress":     0,
			"rdb_last_bgsave_time_sec":   0,
			"connected_clients":          1,
			"blocked_clients":            0,
			"used_memory_dataset":        13546,
			"total_system_memory_human":  "16.00G",
			"maxmemory_policy":           "noeviction",
			"aof_last_cow_size":          0,
			"instantaneous_ops_per_sec":  0,
			"slave_expires_tracked_keys": 0,
			"db0":                    "keys=1,expires=0,avg_ttl=0",
			"role":                   "master",
			"connected_slaves":       0,
			"redis_version":          "4.0.9",
			"sync_partial_ok":        0,
			"migrate_cached_sockets": 0,
			"os":              "Darwin 18.2.0 x86_64",
			"keyspace_misses": "1",
			"hz":              "10",
			"used_memory_lua_human":          "37.00K",
			"mem_allocator":                  "libc",
			"aof_last_bgrewrite_status":      "ok",
			"total_net_input_bytes":          "0",
			"total_net_output_bytes":         "871587",
			"used_memory_peak_human":         "1019.52K",
			"rdb_changes_since_last_save":    "0",
			"total_connections_received":     "323",
			"pubsub_channels":                "0",
			"latest_fork_usec":               "1234",
			"used_memory_rss":                "557056",
			"loading":                        "0",
			"aof_rewrite_in_progress":        "0",
			"rejected_connections":           "0",
			"client_longest_output_list":     "0",
			"used_memory_dataset_perc":       "21.42%",
			"evicted_keys":                   "0",
			"active_defrag_key_hits":         "0",
			"used_cpu_sys_children":          "0.00",
			"gcc_version":                    "4.2.1",
			"executable":                     "/usr/local/opt/redis/bin/redis-server",
			"active_defrag_running":          "0",
			"aof_last_rewrite_time_sec":      "-1",
			"used_cpu_user_children":         "0.00",
			"atomicvar_api":                  "atomic-builtin",
			"rdb_last_cow_size":              "0",
			"aof_enabled":                    "0",
			"pubsub_patterns":                "0",
			"master_repl_offset":             "0",
			"repl_backlog_first_byte_offset": "0",
			"tcp_port":                       "6379",
			"lru_clock":                      "6777537",
			"used_memory_peak":               "1043984",
			"aof_rewrite_scheduled":          "0",
			"expired_stale_perc":             "0.00",
			"repl_backlog_active":            "0",
			"arch_bits":                      "64",
			"run_id":                         "7d1e18e7ccf9484ef3636ff145dbc3ced5fd8f01",
			"rdb_last_save_time":             "1548333171",
			"repl_backlog_histlen":           "0",
			"cluster_enabled":                "0",
			"uptime_in_seconds":              "3349376",
			"second_repl_offset":             "-1",
			"uptime_in_days":                 "38",
			"used_memory_overhead":           "1030438",
			"active_defrag_hits":             "0",
			"used_memory_human":              "1019.52K",
			"used_memory_lua":                "37888",
			"aof_current_rewrite_time_sec":   "-1",
			"sync_partial_err":               "0",
			"repl_backlog_size":              "1048576",
			"config_file":                    "/usr/local/etc/redis.conf",
			"used_memory_startup":            "980736",
			"rdb_current_bgsave_time_sec":    "-1",
			"sync_full":                      "0",
			"used_cpu_sys":                   "422.27",
			"redis_git_sha1":                 "00000000",
			"rdb_last_bgsave_status":         "ok",
			"aof_last_write_status":          "ok",
			"expired_keys":                   "0",
			"redis_build_id":                 "e0c8d37381c486c6",
			"redis_mode":                     "standalone",
			"maxmemory_human":                "0B",
			"zHost":                          "HELLO",
			"used_memory_rss_human":          "544.00K",
			"used_memory_peak_perc":          "100.00%",
			"expired_time_cap_reached_count": "0",
			"active_defrag_misses":           "0",
		},
	}

	config := load.Config{
		Name: "RedisInfo",
		APIs: []load.API{
			{
				Name: "redis",
				Commands: []load.Command{
					{
						Run:     "cat ../../test/payloads/redisInfo.out",
						SplitBy: ":",
					},
					{
						Run:     `echo "zHost:$(echo HELLO)"`,
						SplitBy: ":",
					},
				},
				RemoveKeys:    []string{"human"},
				SnakeToCamel:  true,
				PercToDecimal: true,
				RenameKeys:    map[string]string{"Host": "opSystem"},
				SubParse: []load.Parse{
					{
						Type:    "prefix",
						Key:     "db",
						SplitBy: []string{",", "="},
					},
				},
				CustomAttributes: map[string]string{
					"myCustomAttr": "theValue",
				},
				MetricParser: load.MetricParser{
					Metrics: map[string]string{
						"totalNetInputBytes": "RATE",
						"rate$":              "RATE",
					},
					Namespace: load.Namespace{
						ExistingAttr: []string{"redisVersion", "tcpPort"},
					},
					AutoSet: true,
				},
			},
		},
	}

	dataStore := []interface{}{}
	RunCommands(&config, config.APIs[0], &dataStore)

	for key := range dataStore[0].(map[string]interface{}) {
		if fmt.Sprintf("%v", dataStore[0].(map[string]interface{})[key]) != fmt.Sprintf("%v", dataStoreExpected[0].(map[string]interface{})[key]) {
			t.Errorf(fmt.Sprintf("doesnt match %v : %v - %v", key, dataStore[0].(map[string]interface{})[key], dataStoreExpected[0].(map[string]interface{})[key]))
		}
	}

}

func TestDf(t *testing.T) {
	load.Refresh()
	config := load.Config{
		Name: "dfFlex",
		APIs: []load.API{
			{
				Name:  "df",
				Shell: "/bin/sh",
				Commands: []load.Command{
					{
						Run:   "cat ../../test/payloads/df.out",
						Split: "horizontal",
						SetHeader: []string{
							"fs", "512Blocks", "used", "available", "capacity", "iused", "ifree", "iusedPerc", "mountedOn",
						},
						RegexMatch: false,
						SplitBy:    `\s{1,}`,
						Shell:      "/bin/sh",
					},
				},
			},
		},
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
	RunCommands(&config, config.APIs[0], &dataStore)

	for key := range dataStore[0].(map[string]interface{}) {
		if fmt.Sprintf("%v", dataStore[0].(map[string]interface{})[key]) != fmt.Sprintf("%v", dataStoreExpected[0].(map[string]interface{})[key]) {
			t.Errorf(fmt.Sprintf("doesnt match %v : %v - %v", key, dataStore[0].(map[string]interface{})[key], dataStoreExpected[0].(map[string]interface{})[key]))
		}
	}

}

func TestDf2(t *testing.T) {
	load.Refresh()
	config := load.Config{
		Name: "dfFlex",
		APIs: []load.API{
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
		},
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
	RunCommands(&config, config.APIs[0], &dataStore)

	for key := range dataStore[0].(map[string]interface{}) {
		if fmt.Sprintf("%v", dataStore[0].(map[string]interface{})[key]) != fmt.Sprintf("%v", dataStoreExpected[0].(map[string]interface{})[key]) {
			t.Errorf(fmt.Sprintf("doesnt match %v : %v - %v", key, dataStore[0].(map[string]interface{})[key], dataStoreExpected[0].(map[string]interface{})[key]))
		}
	}
}
