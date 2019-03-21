package parser

// func TestPrometheus(t *testing.T) {
// 	dataStore := []interface{}{}

// 	// create a listener with desired port
// 	l, _ := net.Listen("tcp", "127.0.0.1:9122")
// 	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 		rw.Header().Set("Content-Type", "text/plain; version=0.0.4")
// 		fileData, _ := ioutil.ReadFile("../../test/payloads/prometheusRedis.out")
// 		_, err := rw.Write(fileData)
// 		logger.Flex("debug", err, "failed to write", false)
// 	}))
// 	// NewUnstartedServer creates a listener. Close listener and replace with the one we created.
// 	ts.Listener.Close()
// 	ts.Listener = l
// 	// Start the server.
// 	ts.Start()

// 	config := load.Config{
// 		APIs: []load.API{
// 			{
// 				Name: "prometheusTest",
// 				URL:  "http://localhost:9122",
// 				Prometheus: load.Prometheus{
// 					Enable:         true,
// 					FlattenedEvent: "testPromFlatten",
// 					KeyMerge:       []string{"cmd"},
// 					SampleKeys: map[string]string{
// 						"prometheusRedisDbSample": "db",
// 					},
// 					CustomAttributes: map[string]string{
// 						"abc": "def",
// 					},
// 					Histogram: true,
// 					Summary:   true,
// 				},
// 				RemoveKeys: []string{"go_"},
// 			},
// 		},
// 	}

// 	expectedDatastore := []interface{}{
// 		map[string]interface{}{
// 			"0.000000": "1.7488e-05",
// 			"0.250000": "2.0739e-05",
// 			"0.500000": "8.6961e-05",
// 			"0.750000": "0.000358332",
// 			"1.000000": "0.000606444",
// 			"abc":      "def", "count": "8",
// 			"event_type": "testPromFlatten",
// 			"help":       "A summary of the GC invocation durations.",
// 			"name":       "go_gc_duration_seconds",
// 			"sum":        "0.001253561",
// 			"type":       "SUMMARY",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db1",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db10",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db12",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db3",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db4",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db9",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db15",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db2",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db5",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db6",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                      "def",
// 			"addr":                     "192.168.0.101:6379",
// 			"alias":                    "",
// 			"db":                       "db0",
// 			"event_type":               "prometheusRedisDbSample",
// 			"redis_db_avg_ttl_seconds": "0",
// 			"redis_db_keys":            "1",
// 			"redis_db_keys_expiring":   "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db11",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db13",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db14",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db7",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":                    "def",
// 			"addr":                   "192.168.0.101:6379",
// 			"alias":                  "",
// 			"db":                     "db8",
// 			"event_type":             "prometheusRedisDbSample",
// 			"redis_db_keys":          "0",
// 			"redis_db_keys_expiring": "0",
// 		},
// 		map[string]interface{}{
// 			"abc":        "def",
// 			"event_type": "testPromFlatten",
// 			// "go_gc_duration_seconds.count":                 "8",
// 			// "go_gc_duration_seconds.sum":                   "0.001253561",
// 			// "go_goroutines":                                "7",
// 			// "go_memstats_alloc_bytes":                      "3.128376e+06",
// 			// "go_memstats_alloc_bytes_total":                "2.6926608e+07",
// 			// "go_memstats_buck_hash_sys_bytes":              "1.444544e+06",
// 			// "go_memstats_frees_total":                      "37213",
// 			// "go_memstats_gc_sys_bytes":                     "2.377728e+06",
// 			// "go_memstats_heap_alloc_bytes":                 "3.128376e+06",
// 			// "go_memstats_heap_idle_bytes":                  "6.2398464e+07",
// 			// "go_memstats_heap_inuse_bytes":                 "3.989504e+06",
// 			// "go_memstats_heap_objects":                     "4418",
// 			// "go_memstats_heap_released_bytes_total":        "0",
// 			// "go_memstats_heap_sys_bytes":                   "6.6387968e+07",
// 			// "go_memstats_last_gc_time_seconds":             "1.5523949892526178e+09",
// 			// "go_memstats_lookups_total":                    "0",
// 			// "go_memstats_mallocs_total":                    "41631",
// 			// "go_memstats_mcache_inuse_bytes":               "6944",
// 			// "go_memstats_mcache_sys_bytes":                 "16384",
// 			// "go_memstats_mspan_inuse_bytes":                "27792",
// 			// "go_memstats_mspan_sys_bytes":                  "32768",
// 			// "go_memstats_next_gc_bytes":                    "4.194304e+06",
// 			// "go_memstats_other_sys_bytes":                  "1.306168e+06",
// 			// "go_memstats_stack_inuse_bytes":                "720896",
// 			// "go_memstats_stack_sys_bytes":                  "720896",
// 			// "go_memstats_sys_bytes":                        "7.2286456e+07",
// 			"process_cpu_seconds_total":                    "0.19",
// 			"process_max_fds":                              "1.048576e+06",
// 			"process_open_fds":                             "8",
// 			"process_resident_memory_bytes":                "1.1956224e+07",
// 			"process_start_time_seconds":                   "1.55239479564e+09",
// 			"process_virtual_memory_bytes":                 "1.13074176e+08",
// 			"redis_aof_current_rewrite_duration_sec":       "-1",
// 			"redis_aof_enabled":                            "0",
// 			"redis_aof_last_bgrewrite_status":              "1",
// 			"redis_aof_last_cow_size_bytes":                "0",
// 			"redis_aof_last_rewrite_duration_sec":          "-1",
// 			"redis_aof_last_write_status":                  "1",
// 			"redis_aof_rewrite_in_progress":                "0",
// 			"redis_aof_rewrite_scheduled":                  "0",
// 			"redis_blocked_clients":                        "0",
// 			"redis_client_biggest_input_buf":               "0",
// 			"redis_client_longest_output_list":             "0",
// 			"redis_cluster_enabled":                        "0",
// 			"redis_commands_duration_seconds_total.config": "0.000289",
// 			"redis_commands_processed_total":               "1",
// 			"redis_commands_total.config":                  "1",
// 			"redis_config_maxclients":                      "10000",
// 			"redis_config_maxmemory":                       "0",
// 			"redis_connected_clients":                      "1",
// 			"redis_connected_slaves":                       "0",
// 			"redis_connections_received_total":             "1",
// 			"redis_evicted_keys_total":                     "0",
// 			"redis_expired_keys_total":                     "0",
// 			"redis_exporter_build_info":                    "1",
// 			"redis_exporter_last_scrape_duration_seconds":  "0.007935711",
// 			"redis_exporter_last_scrape_error":             "0",
// 			"redis_exporter_scrapes_total":                 "28",
// 			"redis_instance_info":                          "1",
// 			"redis_instantaneous_input_kbps":               "0",
// 			"redis_instantaneous_ops_per_sec":              "0",
// 			"redis_instantaneous_output_kbps":              "0",
// 			"redis_keyspace_hits_total":                    "0",
// 			"redis_keyspace_misses_total":                  "0",
// 			"redis_latest_fork_usec":                       "0",
// 			"redis_loading_dump_file":                      "0",
// 			"redis_master_repl_offset":                     "0",
// 			"redis_memory_max_bytes":                       "0",
// 			"redis_memory_used_bytes":                      "1.032464e+06",
// 			"redis_memory_used_lua_bytes":                  "37888",
// 			"redis_memory_used_peak_bytes":                 "1.032464e+06",
// 			"redis_memory_used_rss_bytes":                  "2.166784e+06",
// 			"redis_net_input_bytes_total":                  "55",
// 			"redis_net_output_bytes_total":                 "3056",
// 			"redis_process_id":                             "72295",
// 			"redis_pubsub_channels":                        "0",
// 			"redis_pubsub_patterns":                        "0",
// 			"redis_rdb_bgsave_in_progress":                 "0",
// 			"redis_rdb_changes_since_last_save":            "0",
// 			"redis_rdb_current_bgsave_duration_sec":        "-1",
// 			"redis_rdb_last_bgsave_duration_sec":           "-1",
// 			"redis_rdb_last_bgsave_status":                 "1",
// 			"redis_rdb_last_cow_size_bytes":                "0",
// 			"redis_rdb_last_save_timestamp_seconds":        "1.552395036e+09",
// 			"redis_rejected_connections_total":             "0",
// 			"redis_replication_backlog_bytes":              "1.048576e+06",
// 			"redis_slowlog_last_id":                        "0",
// 			"redis_slowlog_length":                         "0",
// 			"redis_start_time_seconds":                     "1.552395036e+09",
// 			"redis_total_system_memory_bytes":              "1.7179869184e+10",
// 			"redis_up":                                     "1",
// 			"redis_uptime_in_seconds":                      "1",
// 			"redis_used_cpu_sys":                           "0.01",
// 			"redis_used_cpu_sys_children":                  "0",
// 			"redis_used_cpu_user":                          "0.01",
// 			"redis_used_cpu_user_children":                 "0",
// 		},
// 	}

// 	doLoop := true
// 	RunHTTP(&doLoop, &config, config.APIs[0], &config.APIs[0].URL, &dataStore)

// 	if len(expectedDatastore) != len(dataStore) {
// 		t.Errorf("Incorrect number of samples generated expected: %d, got: %d", len(expectedDatastore), len(dataStore))
// 		t.Errorf("%v", (dataStore))
// 	}

// 	for _, sample := range expectedDatastore {
// 		switch sample := sample.(type) {
// 		case map[string]interface{}:
// 			for _, rSample := range dataStore {
// 				switch recSample := rSample.(type) {
// 				case map[string]interface{}:

// 					if fmt.Sprintf("%v", recSample["event_type"]) == "prometheusRedisDbSample" && fmt.Sprintf("%v", recSample["db"]) == fmt.Sprintf("%v", sample["db"]) {
// 						for key := range sample {
// 							if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
// 								t.Errorf("dbSample %v want %v, got %v", key, sample[key], recSample[key])
// 							}
// 						}
// 					}

// 					if fmt.Sprintf("%v", recSample["event_type"]) == "testPromFlatten" && fmt.Sprintf("%v", sample["event_type"]) == "testPromFlatten" && fmt.Sprintf("%v", sample["type"]) == fmt.Sprintf("%v", recSample["type"]) {
// 						for key := range sample {
// 							if fmt.Sprintf("%v", sample[key]) != fmt.Sprintf("%v", recSample[key]) {
// 								t.Errorf("dbSample %v want %v, got %v", key, sample[key], recSample[key])
// 							}
// 						}
// 					}
// 				}

// 			}
// 		}

// 	}
// }
