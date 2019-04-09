package processor

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// testSamples as samples could be generated in different orders, so we test per sample
func testSamples(expectedSamples []string, entityMetrics []*metric.Set, t *testing.T) {
	if len(entityMetrics) != len(expectedSamples) {
		t.Errorf("Missing samples, got: %v, want: %v.", (entityMetrics), (expectedSamples))

		// t.Errorf("Missing samples, got: %v, want: %v.", len(entityMetrics), len(expectedSamples))
	}
	for _, expectedSample := range expectedSamples {
		matchedSample := false
		for _, sample := range entityMetrics {
			out, err := sample.MarshalJSON()
			if err != nil {
				logger.Flex("debug", err, "failed to marshal", false)
			} else {
				if expectedSample == string(out) {
					matchedSample = true
					break
				}
			}
		}
		if !matchedSample {
			completeMetrics, _ := json.Marshal(entityMetrics)
			t.Errorf("Unable to find expected payload, received: %v, want: %v.", string(completeMetrics), expectedSample)
		}
	}
}

func TestConfigDir(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestReadJsonCmdDir", "nri-flex")
	load.Args.ConfigDir = "../../test/configs/"

	var ymls []load.Config
	var files []os.FileInfo

	path := filepath.FromSlash(load.Args.ConfigDir)
	var err error
	files, err = ioutil.ReadDir(path)
	logger.Flex("fatal", err, "failed to read config dir: "+load.Args.ConfigDir, false)

	LoadConfigFiles(&ymls, files, path) // load standard configs if available
	RunConfigFiles(&ymls)

	expectedSamples := []string{
		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.newrelic.nri-flex",` +
			`"integration_version":"Unknown-SNAPSHOT","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

// func TestConfigFile(t *testing.T) {
// 	load.Refresh()
// 	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
// 	load.Entity, _ = i.Entity("TestReadJsonCmd", "nri-flex")
// 	load.Args.ConfigFile = "../../test/configs/json-read-cmd-example.yml"
// 	runIntegration(i)
// 	expectedSamples := []string{
// 		`{"configsProcessed":1,"eventCount":1,"eventDropCount":0,"event_type":"flexStatusSample"}`,
// 		`{"completed":"false","event_type":"commandJsonOutSample","id":1,"integration_name":"com.newrelic.nri-flex",` +
// 			`"integration_version":"Unknown-SNAPSHOT","myCustomAttr":"theValue","title":"delectus aut autem","userId":1}`}
// 	testSamples(expectedSamples, load.Entity.Metrics, t)
// }

func TestSubEnvVariables(t *testing.T) {
	str := " hi there $$PWD bye"
	SubEnvVariables(&str)
	if strings.Count(str, "$$") != 0 {
		t.Errorf("failed to sub all variables %v", str)
	}
}

func TestE2E_Vault(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestVault", "nri-flex")
	expectedSample :=
		`{"baseUrl":"http://127.0.0.1:32768/v1/","event_type":"vaultStatus","health.api.StatusCode":200,` +
			`"health.cluster_id":"f0759ac2-22a4-6432-81b1-af3f599d107a","health.cluster_name":"vault-cluster-3abb3c80",` +
			`"health.initialized":"true","health.performance_standby":"false","health.replication_dr_mode":"disabled",` +
			`"health.replication_performance_mode":"disabled","health.sealed":"false","health.server_time_utc":1550325568,` +
			`"health.standby":"false","health.version":"1.0.2","integration_name":"com.newrelic.nri-flex",` +
			`"integration_version":"Unknown-SNAPSHOT","key.api.StatusCode":200,"key.auth":"\u003cnil\u003e","key.data.install_time":"2019-02-16T13:44:21.716366Z",` +
			`"key.data.term":1,"key.install_time":"2019-02-16T13:44:21.716366Z","key.lease_duration":0,"key.lease_id":"",` +
			`"key.renewable":"false","key.request_id":"232fa0aa-d7f6-6ac4-c8af-b657b3642ef5","key.term":1,"key.warnings":"\u003cnil\u003e",` +
			`"key.wrap_info":"\u003cnil\u003e","ldr.api.StatusCode":200,"ldr.ha_enabled":"false","ldr.is_self":"false","ldr.leader_address":"",` +
			`"ldr.leader_cluster_address":"","ldr.performance_standby":"false","ldr.performance_standby_last_remote_wal":0,"myVaultNode":"NodeABC",` +
			`"seal.api.StatusCode":200,"seal.cluster_id":"f0759ac2-22a4-6432-81b1-af3f599d107a","seal.cluster_name":"vault-cluster-3abb3c80",` +
			`"seal.initialized":"true","seal.migration":"false","seal.n":1,"seal.nonce":"","seal.progress":0,"seal.recovery_seal":"false",` +
			`"seal.sealed":"false","seal.t":1,"seal.type":"shamir","seal.version":"1.0.2"}`

	var x interface{}
	err := json.Unmarshal([]byte(expectedSample), &x)

	if err != nil {
		t.Errorf("Failed to unmarshal %v", err.Error())
	}

	expectedSampleData := x.(map[string]interface{})

	config := load.Config{
		Name: "vaultFlex",
		Global: load.Global{
			BaseURL: "http://127.0.0.1:32768/v1/",
			Headers: map[string]string{
				"X-Vault-Token": "myroot",
			},
		},
		CustomAttributes: map[string]string{
			"myVaultNode": "NodeABC",
		},
		APIs: []load.API{
			load.API{
				EventType: "vaultKeyStatus",
				File:      "../../test/payloads/vaultKeyStatus.json",
				Prefix:    "key.",
				Merge:     "vaultStatus",
			},
			load.API{
				EventType: "vaultHealthStatus",
				File:      "../../test/payloads/vaultHealthStatus.json",
				Prefix:    "health.",
				Merge:     "vaultStatus",
			},
			load.API{
				EventType: "vaultLeaderStatus",
				File:      "../../test/payloads/vaultLeaderStatus.json",
				Prefix:    "ldr.",
				Merge:     "vaultStatus",
			},
			load.API{
				EventType: "vaultSealStatus",
				File:      "../../test/payloads/vaultSealStatus.json",
				Prefix:    "seal.",
				Merge:     "vaultStatus",
			},
		},
	}

	RunConfig(config)

	if len(load.Entity.Metrics) > 1 {
		t.Errorf("Too many samples expected: 1, got: %v", len(load.Entity.Metrics))
	}

	for key := range expectedSampleData {
		if load.Entity.Metrics[0].Metrics[key] == nil && !strings.Contains(key, "api.StatusCode") {
			t.Errorf("Missing attribute %v", key)
		} else if !strings.Contains(key, "api.StatusCode") && load.Entity.Metrics[0].Metrics[key] != nil {
			if load.Entity.Metrics[0].Metrics[key] != expectedSampleData[key] {
				t.Errorf("Data mismatch received: %v: %v, want: %v.", key, load.Entity.Metrics[0].Metrics[key], expectedSampleData[key])
			}
		}
	}

}

func TestE2E_ElasticSearch(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestEtcd", "nri-flex")
	expectedSamples := []string{
		`{"baseUrl":"http://localhost:9200/","completion.size_in_bytes":0,"docs.count":0,"docs.deleted":0,"event_type":"elasticsearchTotalSample","fielddata.evictions":0,"fielddata.memory_size_in_bytes":0,"flush.periodic":0,"flush.total":10,"flush.total_time_in_millis":5,"get.current":0,"get.exists_time_in_millis":0,"get.exists_total":0,"get.missing_time_in_millis":0,"get.missing_total":0,"get.time_in_millis":0,"get.total":0,"indexing.delete_current":0,"indexing.delete_time_in_millis":0,"indexing.delete_total":0,"indexing.index_current":0,"indexing.index_failed":0,"indexing.index_time_in_millis":0,"indexing.index_total":0,"indexing.is_throttled":"false","indexing.noop_update_total":0,"indexing.throttle_time_in_millis":0,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","merges.current":0,"merges.current_docs":0,"merges.current_size_in_bytes":0,"merges.total":0,"merges.total_auto_throttle_in_bytes":209715200,"merges.total_docs":0,"merges.total_size_in_bytes":0,"merges.total_stopped_time_in_millis":0,"merges.total_throttled_time_in_millis":0,"merges.total_time_in_millis":0,"query_cache.cache_count":0,"query_cache.cache_size":0,"query_cache.evictions":0,"query_cache.hit_count":0,"query_cache.memory_size_in_bytes":0,"query_cache.miss_count":0,"query_cache.total_count":0,"recovery.current_as_source":0,"recovery.current_as_target":0,"recovery.throttle_time_in_millis":0,"refresh.listeners":0,"refresh.total":50,"refresh.total_time_in_millis":0,"request_cache.evictions":0,"request_cache.hit_count":0,"request_cache.memory_size_in_bytes":0,"request_cache.miss_count":0,"search.fetch_current":0,"search.fetch_time_in_millis":0,"search.fetch_total":0,"search.open_contexts":0,"search.query_current":0,"search.query_time_in_millis":0,"search.query_total":0,"search.scroll_current":0,"search.scroll_time_in_millis":0,"search.scroll_total":0,"search.suggest_current":0,"search.suggest_time_in_millis":0,"search.suggest_total":0,"segments.count":0,"segments.doc_values_memory_in_bytes":0,"segments.fixed_bit_set_memory_in_bytes":0,"segments.index_writer_memory_in_bytes":0,"segments.max_unsafe_auto_id_timestamp":-1,"segments.memory_in_bytes":0,"segments.norms_memory_in_bytes":0,"segments.points_memory_in_bytes":0,"segments.stored_fields_memory_in_bytes":0,"segments.term_vectors_memory_in_bytes":0,"segments.terms_memory_in_bytes":0,"segments.version_map_memory_in_bytes":0,"store.size_in_bytes":2610,"translog.earliest_last_modified_age":0,"translog.operations":0,"translog.size_in_bytes":1100,"translog.uncommitted_operations":0,"translog.uncommitted_size_in_bytes":550,"warmer.current":0,"warmer.total":20,"warmer.total_time_in_millis":1}`,
		`{"baseUrl":"http://localhost:9200/","completion.size_in_bytes":0,"docs.count":0,"docs.deleted":0,"event_type":"elasticsearchPrimarySample","fielddata.evictions":0,"fielddata.memory_size_in_bytes":0,"flush.periodic":0,"flush.total":10,"flush.total_time_in_millis":5,"get.current":0,"get.exists_time_in_millis":0,"get.exists_total":0,"get.missing_time_in_millis":0,"get.missing_total":0,"get.time_in_millis":0,"get.total":0,"indexing.delete_current":0,"indexing.delete_time_in_millis":0,"indexing.delete_total":0,"indexing.index_current":0,"indexing.index_failed":0,"indexing.index_time_in_millis":0,"indexing.index_total":0,"indexing.is_throttled":"false","indexing.noop_update_total":0,"indexing.throttle_time_in_millis":0,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","merges.current":0,"merges.current_docs":0,"merges.current_size_in_bytes":0,"merges.total":0,"merges.total_auto_throttle_in_bytes":209715200,"merges.total_docs":0,"merges.total_size_in_bytes":0,"merges.total_stopped_time_in_millis":0,"merges.total_throttled_time_in_millis":0,"merges.total_time_in_millis":0,"query_cache.cache_count":0,"query_cache.cache_size":0,"query_cache.evictions":0,"query_cache.hit_count":0,"query_cache.memory_size_in_bytes":0,"query_cache.miss_count":0,"query_cache.total_count":0,"recovery.current_as_source":0,"recovery.current_as_target":0,"recovery.throttle_time_in_millis":0,"refresh.listeners":0,"refresh.total":50,"refresh.total_time_in_millis":0,"request_cache.evictions":0,"request_cache.hit_count":0,"request_cache.memory_size_in_bytes":0,"request_cache.miss_count":0,"search.fetch_current":0,"search.fetch_time_in_millis":0,"search.fetch_total":0,"search.open_contexts":0,"search.query_current":0,"search.query_time_in_millis":0,"search.query_total":0,"search.scroll_current":0,"search.scroll_time_in_millis":0,"search.scroll_total":0,"search.suggest_current":0,"search.suggest_time_in_millis":0,"search.suggest_total":0,"segments.count":0,"segments.doc_values_memory_in_bytes":0,"segments.fixed_bit_set_memory_in_bytes":0,"segments.index_writer_memory_in_bytes":0,"segments.max_unsafe_auto_id_timestamp":-1,"segments.memory_in_bytes":0,"segments.norms_memory_in_bytes":0,"segments.points_memory_in_bytes":0,"segments.stored_fields_memory_in_bytes":0,"segments.term_vectors_memory_in_bytes":0,"segments.terms_memory_in_bytes":0,"segments.version_map_memory_in_bytes":0,"store.size_in_bytes":2610,"translog.earliest_last_modified_age":0,"translog.operations":0,"translog.size_in_bytes":1100,"translog.uncommitted_operations":0,"translog.uncommitted_size_in_bytes":550,"warmer.current":0,"warmer.total":20,"warmer.total_time_in_millis":1}`,
		`{"baseUrl":"http://localhost:9200/","event_type":"elasticsearchShardSample","failed":0,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","successful":10,"total":20}`,
		`{"active_primary_shards":10,"active_shards":10,"active_shards_percent_as_number":50,"baseUrl":"http://localhost:9200/","cluster_name":"docker-cluster","delayed_unassigned_shards":0,"event_type":"elasticsearchClusterHealthSample","initializing_shards":0,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","number_of_data_nodes":1,"number_of_in_flight_fetch":0,"number_of_nodes":1,"number_of_pending_tasks":0,"relocating_shards":0,"status":"yellow","task_max_waiting_in_queue_millis":0,"timed_out":"false","unassigned_shards":10}`,
		// `{"baseUrl":"http://localhost:9200/","event_type":"elasticsearchPendingClusterTaskSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"}`,
		`{"attributes.ml.enabled":"true","attributes.ml.machine_memory":2095869952,"attributes.ml.max_open_jobs":20,"attributes.xpack.installed":"true",` +
			`"baseUrl":"http://localhost:9200/","breakers.accounting.estimated_size":"0b","breakers.accounting.estimated_size_in_bytes":0,"breakers.accounting.limit_size":"990.7mb",` +
			`"breakers.accounting.limit_size_in_bytes":1038876672,"breakers.accounting.overhead":1,"breakers.accounting.tripped":0,"breakers.fielddata.estimated_size":"0b","breakers.fielddata.estimated_size_in_bytes":0,"breakers.fielddata.limit_size":"594.4mb",` +
			`"breakers.fielddata.limit_size_in_bytes":623326003,"breakers.fielddata.overhead":1.03,"breakers.fielddata.tripped":0,"breakers.in_flight_requests.estimated_size":"0b","breakers.in_flight_requests.estimated_size_in_bytes":0,` +
			`"breakers.in_flight_requests.limit_size":"990.7mb","breakers.in_flight_requests.limit_size_in_bytes":1038876672,"breakers.in_flight_requests.overhead":1,"breakers.in_flight_requests.tripped":0,"breakers.parent.estimated_size":"0b",` +
			`"breakers.parent.estimated_size_in_bytes":0,"breakers.parent.limit_size":"693.5mb","breakers.parent.limit_size_in_bytes":727213670,"breakers.parent.overhead":1,"breakers.parent.tripped":0,"breakers.request.estimated_size":"0b",` +
			`"breakers.request.estimated_size_in_bytes":0,"breakers.request.limit_size":"594.4mb","breakers.request.limit_size_in_bytes":623326003,"breakers.request.overhead":1,"breakers.request.tripped":0,"cluster_name":"docker-cluster",` +
			`"event_type":"elasticsearchNodeSample","host":"172.24.0.2","http.current_open":2,"http.total_opened":19,"indices.completion.size_in_bytes":0,"indices.docs.count":0,"indices.docs.deleted":0,"indices.fielddata.evictions":0,` +
			`"indices.fielddata.memory_size_in_bytes":0,"indices.flush.periodic":0,"indices.flush.total":10,"indices.flush.total_time_in_millis":5,"indices.get.current":0,"indices.get.exists_time_in_millis":0,"indices.get.exists_total":0,` +
			`"indices.get.missing_time_in_millis":0,"indices.get.missing_total":0,"indices.get.time_in_millis":0,"indices.get.total":0,"indices.indexing.delete_current":0,"indices.indexing.delete_time_in_millis":0,"indices.indexing.delete_total":0,` +
			`"indices.indexing.index_current":0,"indices.indexing.index_failed":0,"indices.indexing.index_time_in_millis":0,"indices.indexing.index_total":0,"indices.indexing.is_throttled":"false","indices.indexing.noop_update_total":0,` +
			`"indices.indexing.throttle_time_in_millis":0,"indices.merges.current":0,"indices.merges.current_docs":0,"indices.merges.current_size_in_bytes":0,"indices.merges.total":0,"indices.merges.total_auto_throttle_in_bytes":209715200,` +
			`"indices.merges.total_docs":0,"indices.merges.total_size_in_bytes":0,"indices.merges.total_stopped_time_in_millis":0,"indices.merges.total_throttled_time_in_millis":0,"indices.merges.total_time_in_millis":0,"indices.query_cache.cache_count":0,` +
			`"indices.query_cache.cache_size":0,"indices.query_cache.evictions":0,"indices.query_cache.hit_count":0,"indices.query_cache.memory_size_in_bytes":0,"indices.query_cache.miss_count":0,"indices.query_cache.total_count":0,"indices.recovery.current_as_source":0,` +
			`"indices.recovery.current_as_target":0,"indices.recovery.throttle_time_in_millis":0,"indices.refresh.listeners":0,"indices.refresh.total":50,"indices.refresh.total_time_in_millis":0,"indices.request_cache.evictions":0,"indices.request_cache.hit_count":0,` +
			`"indices.request_cache.memory_size_in_bytes":0,"indices.request_cache.miss_count":0,"indices.search.fetch_current":0,"indices.search.fetch_time_in_millis":0,"indices.search.fetch_total":0,"indices.search.open_contexts":0,"indices.search.query_current":0,` +
			`"indices.search.query_time_in_millis":0,"indices.search.query_total":0,"indices.search.scroll_current":0,"indices.search.scroll_time_in_millis":0,"indices.search.scroll_total":0,"indices.search.suggest_current":0,"indices.search.suggest_time_in_millis":0,` +
			`"indices.search.suggest_total":0,"indices.segments.count":0,"indices.segments.doc_values_memory_in_bytes":0,"indices.segments.fixed_bit_set_memory_in_bytes":0,"indices.segments.index_writer_memory_in_bytes":0,` +
			`"indices.segments.max_unsafe_auto_id_timestamp":-1,"indices.segments.memory_in_bytes":0,"indices.segments.norms_memory_in_bytes":0,"indices.segments.points_memory_in_bytes":0,"indices.segments.stored_fields_memory_in_bytes":0,` +
			`"indices.segments.term_vectors_memory_in_bytes":0,"indices.segments.terms_memory_in_bytes":0,"indices.segments.version_map_memory_in_bytes":0,"indices.store.size_in_bytes":2610,"indices.translog.earliest_last_modified_age":0,` +
			`"indices.translog.operations":0,"indices.translog.size_in_bytes":1100,"indices.translog.uncommitted_operations":0,"indices.translog.uncommitted_size_in_bytes":550,"indices.warmer.current":0,"indices.warmer.total":20,` +
			`"indices.warmer.total_time_in_millis":1,"ingest.total.count":0,"ingest.total.current":0,"ingest.total.failed":0,"ingest.total.time_in_millis":0,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","ip":"172.24.0.2:9300","jvm.buffer_pools.direct.count":29,"jvm.buffer_pools.direct.total_capacity_in_bytes":134781616,"jvm.buffer_pools.direct.used_in_bytes":134781617,"jvm.buffer_pools.mapped.count":0,"jvm.buffer_pools.mapped.total_capacity_in_bytes":0,"jvm.buffer_pools.mapped.used_in_bytes":0,"jvm.classes.current_loaded_count":15250,"jvm.classes.total_loaded_count":15250,"jvm.classes.total_unloaded_count":0,"jvm.gc.collectors.old.collection_count":2,"jvm.gc.collectors.old.collection_time_in_millis":90,"jvm.gc.collectors.young.collection_count":7,"jvm.gc.collectors.young.collection_time_in_millis":216,"jvm.mem.heap_committed_in_bytes":1038876672,"jvm.mem.heap_max_in_bytes":1038876672,"jvm.mem.heap_used_in_bytes":265996568,"jvm.mem.heap_used_percent":25,"jvm.mem.non_heap_committed_in_bytes":117010432,"jvm.mem.non_heap_used_in_bytes":108058336,"jvm.mem.pools.old.max_in_bytes":724828160,"jvm.mem.pools.old.peak_max_in_bytes":724828160,"jvm.mem.pools.old.peak_used_in_bytes":183421352,"jvm.mem.pools.old.used_in_bytes":183421352,"jvm.mem.pools.survivor.max_in_bytes":34865152,"jvm.mem.pools.survivor.peak_max_in_bytes":34865152,"jvm.mem.pools.survivor.peak_used_in_bytes":34865152,"jvm.mem.pools.survivor.used_in_bytes":29387600,"jvm.mem.pools.young.max_in_bytes":279183360,"jvm.mem.pools.young.peak_max_in_bytes":279183360,"jvm.mem.pools.young.peak_used_in_bytes":279183360,"jvm.mem.pools.young.used_in_bytes":53187616,"jvm.threads.count":34,"jvm.threads.peak_count":34,"jvm.timestamp":1550488706122,"jvm.uptime_in_millis":561366,"name":"Ui20L36","node.id":"Ui20L36kQle5ZEHCpWUvrw","os.cgroup.cpu.control_group":"/","os.cgroup.cpu.stat.number_of_elapsed_periods":0,"os.cgroup.cpu.stat.number_of_times_throttled":0,"os.cgroup.cpu.stat.time_throttled_nanos":0,"os.cgroup.cpuacct.control_group":"/","os.cgroup.cpuacct.usage_nanos":38539936769,"os.cgroup.memory.control_group":"/","os.cgroup.memory.limit_in_bytes":9223372036854772000,"os.cgroup.memory.usage_in_bytes":1545510912,"os.cpu.load_average.15m":0.03,"os.cpu.load_average.1m":0.09,"os.cpu.load_average.5m":0.06,"os.cpu.percent":0,"os.mem.free_in_bytes":93130752,"os.mem.free_percent":4,"os.mem.total_in_bytes":2095869952,"os.mem.used_in_bytes":2002739200,"os.mem.used_percent":96,"os.swap.free_in_bytes":1063653376,"os.swap.total_in_bytes":1073737728,"os.swap.used_in_bytes":10084352,"os.timestamp":1550488706122,"parentNodes.failed":0,"parentNodes.successful":1,"parentNodes.total":1,"process.cpu.percent":0,"process.cpu.total_in_millis":38010,"process.max_file_descriptors":1048576,"process.mem.total_virtual_in_bytes":5001908224,"process.open_file_descriptors":271,"process.timestamp":1550488706122,"script.cache_evictions":0,"script.compilations":4,"thread_pool.analyze.active":0,"thread_pool.analyze.completed":0,"thread_pool.analyze.largest":0,"thread_pool.analyze.queue":0,"thread_pool.analyze.rejected":0,"thread_pool.analyze.threads":0,"thread_pool.ccr.active":0,"thread_pool.ccr.completed":0,"thread_pool.ccr.largest":0,"thread_pool.ccr.queue":0,"thread_pool.ccr.rejected":0,"thread_pool.ccr.threads":0,"thread_pool.fetch_shard_started.active":0,"thread_pool.fetch_shard_started.completed":0,"thread_pool.fetch_shard_started.largest":0,"thread_pool.fetch_shard_started.queue":0,"thread_pool.fetch_shard_started.rejected":0,"thread_pool.fetch_shard_started.threads":0,"thread_pool.fetch_shard_store.active":0,"thread_pool.fetch_shard_store.completed":0,"thread_pool.fetch_shard_store.largest":0,"thread_pool.fetch_shard_store.queue":0,"thread_pool.fetch_shard_store.rejected":0,"thread_pool.fetch_shard_store.threads":0,"thread_pool.flush.active":0,"thread_pool.flush.completed":20,"thread_pool.flush.largest":2,"thread_pool.flush.queue":0,"thread_pool.flush.rejected":0,"thread_pool.flush.threads":2,"thread_pool.force_merge.active":0,"thread_pool.force_merge.completed":0,"thread_pool.force_merge.largest":0,"thread_pool.force_merge.queue":0,"thread_pool.force_merge.rejected":0,"thread_pool.force_merge.threads":0,"thread_pool.generic.active":0,"thread_pool.generic.completed":1202,"thread_pool.generic.largest":4,"thread_pool.generic.queue":0,"thread_pool.generic.rejected":0,"thread_pool.generic.threads":4,"thread_pool.get.active":0,"thread_pool.get.completed":0,"thread_pool.get.largest":0,"thread_pool.get.queue":0,"thread_pool.get.rejected":0,"thread_pool.get.threads":0,"thread_pool.index.active":0,"thread_pool.index.completed":0,"thread_pool.index.largest":0,"thread_pool.index.queue":0,"thread_pool.index.rejected":0,"thread_pool.index.threads":0,"thread_pool.listener.active":0,"thread_pool.listener.completed":0,"thread_pool.listener.largest":0,"thread_pool.listener.queue":0,"thread_pool.listener.rejected":0,"thread_pool.listener.threads":0,"thread_pool.management.active":1,"thread_pool.management.completed":81,"thread_pool.management.largest":3,"thread_pool.management.queue":0,"thread_pool.management.rejected":0,"thread_pool.management.threads":3,"thread_pool.ml_autodetect.active":0,"thread_pool.ml_autodetect.completed":0,"thread_pool.ml_autodetect.largest":0,"thread_pool.ml_autodetect.queue":0,"thread_pool.ml_autodetect.rejected":0,"thread_pool.ml_autodetect.threads":0,"thread_pool.ml_datafeed.active":0,"thread_pool.ml_datafeed.completed":0,"thread_pool.ml_datafeed.largest":0,"thread_pool.ml_datafeed.queue":0,"thread_pool.ml_datafeed.rejected":0,"thread_pool.ml_datafeed.threads":0,"thread_pool.ml_utility.active":0,"thread_pool.ml_utility.completed":0,"thread_pool.ml_utility.largest":0,"thread_pool.ml_utility.queue":0,"thread_pool.ml_utility.rejected":0,"thread_pool.ml_utility.threads":0,"thread_pool.refresh.active":0,"thread_pool.refresh.completed":879,"thread_pool.refresh.largest":1,"thread_pool.refresh.queue":0,"thread_pool.refresh.rejected":0,"thread_pool.refresh.threads":1,"thread_pool.rollup_indexing.active":0,"thread_pool.rollup_indexing.completed":0,"thread_pool.rollup_indexing.largest":0,"thread_pool.rollup_indexing.queue":0,"thread_pool.rollup_indexing.rejected":0,"thread_pool.rollup_indexing.threads":0,"thread_pool.search.active":0,"thread_pool.search.completed":0,"thread_pool.search.largest":0,"thread_pool.search.queue":0,"thread_pool.search.rejected":0,"thread_pool.search.threads":0,"thread_pool.search_throttled.active":0,"thread_pool.search_throttled.completed":0,"thread_pool.search_throttled.largest":0,"thread_pool.search_throttled.queue":0,"thread_pool.search_throttled.rejected":0,"thread_pool.search_throttled.threads":0,"thread_pool.security-token-key.active":0,"thread_pool.security-token-key.completed":0,"thread_pool.security-token-key.largest":0,"thread_pool.security-token-key.queue":0,"thread_pool.security-token-key.rejected":0,"thread_pool.security-token-key.threads":0,"thread_pool.snapshot.active":0,"thread_pool.snapshot.completed":0,"thread_pool.snapshot.largest":0,"thread_pool.snapshot.queue":0,"thread_pool.snapshot.rejected":0,"thread_pool.snapshot.threads":0,"thread_pool.warmer.active":0,"thread_pool.warmer.completed":0,"thread_pool.warmer.largest":0,"thread_pool.warmer.queue":0,"thread_pool.warmer.rejected":0,"thread_pool.warmer.threads":0,"thread_pool.watcher.active":0,"thread_pool.watcher.completed":0,"thread_pool.watcher.largest":0,"thread_pool.watcher.queue":0,"thread_pool.watcher.rejected":0,"thread_pool.watcher.threads":0,"thread_pool.write.active":0,"thread_pool.write.completed":0,"thread_pool.write.largest":0,"thread_pool.write.queue":0,"thread_pool.write.rejected":0,"thread_pool.write.threads":0,"timestamp":1550488706118,"transport.rx_count":0,"transport.rx_size_in_bytes":0,"transport.server_open":0,"transport.tx_count":0,"transport.tx_size_in_bytes":0,"transport_address":"172.24.0.2:9300"}`,
	}
	config := load.Config{
		Name: "Prom Range Test",
		Global: load.Global{
			BaseURL: "http://localhost:9200/",
		},
		APIs: []load.API{
			load.API{
				EventType: "elasticsearchTotalSample",
				File:      "../../test/payloads/elasticsearchStats.json",
				StartKey:  []string{"_all", "total"},
			},
			load.API{
				EventType: "elasticsearchPrimarySample",
				Cache:     "../../test/payloads/elasticsearchStats.json",
				StartKey:  []string{"_all", "total"},
			},
			load.API{
				EventType: "elasticsearchShardSample",
				Cache:     "../../test/payloads/elasticsearchStats.json",
				StartKey:  []string{"_shards"},
			},
			load.API{
				EventType: "elasticsearchClusterHealthSample",
				File:      "../../test/payloads/elasticsearchClusterHealth.json",
			},
			// load.API{
			// 	EventType: "elasticsearchClusterPendingTaskSample",
			// 	File:      "../../test/payloads/elasticsearchClusterPendingTask.json",
			// },
			load.API{
				EventType: "elasticsearchNodeSample",
				File:      "../../test/payloads/elasticsearchNodeStats.json",
				SampleKeys: map[string]string{
					"nodes": "nodes>node.id",
				},
				RenameKeys: map[string]string{
					"_nodes": "parentNodes",
				},
				RemoveKeys: []string{"ingest.pipelines.xpack", "roleSampleSamples", "fs."},
			},
		},
	}

	RunConfig(config)
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestRedis(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestRedis", "nri-flex")
	expectedSamples := []string{
		`{"activeDefragHits":0,"activeDefragKeyHits":0,"activeDefragKeyMisses":0,"activeDefragMisses":0,"activeDefragRunning":0,` +
			`"aofCurrentRewriteTimeSec":-1,"aofEnabled":0,"aofLastBgrewriteStatus":"ok","aofLastCowSize":0,"aofLastRewriteTimeSec":-1,` +
			`"aofLastWriteStatus":"ok","aofRewriteInProgress":0,"aofRewriteScheduled":0,"archBits":64,"atomicvarApi":"atomic-builtin",` +
			`"blockedClients":0,"clientBiggestInputBuf":0,"clientLongestOutputList":0,"clusterEnabled":0,"configFile":"/usr/local/etc/redis.conf",` +
			`"connectedClients":1,"connectedSlaves":0,"db0":"keys=1,expires=0,avg_ttl=0","db0.avg_ttl":0,"db0.expires":0,"db0.keys":1,` +
			`"event_type":"redisSample","evictedKeys":0,"executable":"/usr/local/opt/redis/bin/redis-server","expiredKeys":0,"expiredStalePerc":0,` +
			`"expiredTimeCapReachedCount":0,"gccVersion":"4.2.1","hz":10,"instantaneousInputKbps":0,"instantaneousOpsPerSec":0,` +
			`"instantaneousOutputKbps":0,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","keyspaceHits":1,` +
			`"keyspaceMisses":1,"latestForkUsec":1234,"lazyfreePendingObjects":0,"loading":0,"lruClock":6777537,"masterReplOffset":0,` +
			`"masterReplid":"65a558faf348fb6f0cf75b82e4ae6b1cc4128baf","masterReplid2":0,"maxmemory":0,"maxmemoryPolicy":"noeviction",` +
			`"memAllocator":"libc","memFragmentationRatio":0.53,"migrateCachedSockets":0,"multiplexingApi":"kqueue","myCustomAttr":"theValue",` +
			`"namespace":"4.0.9-6379","os":"Darwin 18.2.0 x86_64","processId":1201,"pubsubChannels":0,"pubsubPatterns":0,"rdbBgsaveInProgress":0,` +
			`"rdbChangesSinceLastSave":0,"rdbCurrentBgsaveTimeSec":-1,"rdbLastBgsaveStatus":"ok","rdbLastBgsaveTimeSec":0,"rdbLastCowSize":0,` +
			`"rdbLastSaveTime":1548333171,"redisBuildId":"e0c8d37381c486c6","redisGitDirty":0,"redisGitSha1":0,"redisMode":"standalone",` +
			`"redisVersion":"4.0.9","rejectedConnections":0,"replBacklogActive":0,"replBacklogFirstByteOffset":0,"replBacklogHistlen":0,` +
			`"replBacklogSize":1048576,"role":"master","runId":"7d1e18e7ccf9484ef3636ff145dbc3ced5fd8f01","secondReplOffset":-1,` +
			`"slaveExpiresTrackedKeys":0,"syncFull":0,"syncPartialErr":0,"syncPartialOk":0,"tcpPort":6379,"totalCommandsProcessed":331,` +
			`"totalConnectionsReceived":323,"totalNetInputBytes":0,"totalNetOutputBytes":871587,"totalSystemMemory":17179869184,"uptimeInDays":38,` +
			`"uptimeInSeconds":3349376,"usedCpuSys":422.27,"usedCpuSysChildren":0,"usedCpuUser":200.05,"usedCpuUserChildren":0,"usedMemory":1043984,` +
			`"usedMemoryDataset":13546,"usedMemoryDatasetPerc":21.42,"usedMemoryLua":37888,"usedMemoryOverhead":1030438,"usedMemoryPeak":1043984,` +
			`"usedMemoryPeakPerc":100,"usedMemoryRss":557056,"usedMemoryStartup":980736,"zopSystem":"HELLO"}`,
	}

	config := load.Config{
		Name: "RedisInfo",
		APIs: []load.API{
			load.API{
				Name: "redis",
				Commands: []load.Command{
					load.Command{
						Run:     "cat ../../test/payloads/redisInfo.out",
						SplitBy: ":",
					},
					load.Command{
						Run:     `echo "zHost:$(echo HELLO)"`,
						SplitBy: ":",
					},
				},
				RemoveKeys:    []string{"human"},
				SnakeToCamel:  true,
				PercToDecimal: true,
				RenameKeys:    map[string]string{"Host": "opSystem"},
				SubParse: []load.Parse{
					load.Parse{
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

	RunConfig(config)
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestE2E_PrometheusQueryMetrics(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestPrometheusQueryMetrics", "nri-flex")
	expectedSamples := []string{
		`{"api.StatusCode":200,"baseUrl":"http://localhost:9090/api/v1","data.resultType":"vector","event_type":"promTestQuerySample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","metric.__name__":"up","metric.instance":"localhost:9090","metric.job":"prometheus","status":"success","timestamp":1435781451,"value":1}`,
		`{"api.StatusCode":200,"baseUrl":"http://localhost:9090/api/v1","data.resultType":"vector","event_type":"promTestQuerySample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","metric.__name__":"up","metric.instance":"localhost:9100","metric.job":"node","status":"success","timestamp":1435781451,"value":0}`,
	}

	// create a listener with desired port
	l, _ := net.Listen("tcp", "127.0.0.1:9090")
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		fileData, _ := ioutil.ReadFile("../../test/payloads/pquery.json")
		_, err := rw.Write(fileData)
		logger.Flex("debug", err, "", false)
	})) // NewUnstartedServer creates a listener.

	ts.Listener.Close() // Close listener and replace with the one we created.
	ts.Listener = l
	// Start the server.
	ts.Start()

	config := load.Config{
		Name: "PromRangeFlex",
		Global: load.Global{
			BaseURL: "http://localhost:9090/api/v1",
		},
		APIs: []load.API{
			load.API{
				EventType: "promTestQuerySample",
				URL:       "/",
				// File: "../../test/payloads/pquery.json",
			},
		},
	}

	RunConfig(config)
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestE2E_PrometheusQueryRangeMetrics(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestPrometheusQueryRangeMetrics", "nri-flex")
	expectedSamples := []string{
		`{"baseUrl":"http://localhost:9090/api/v1/","data.resultType":"matrix","event_type":"promTestQueryRangeSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","metric.__name__":"up","metric.instance":"localhost:9090","metric.job":"prometheus","status":"success","timestamp":1435781430,"value":1}`,
		`{"baseUrl":"http://localhost:9090/api/v1/","data.resultType":"matrix","event_type":"promTestQueryRangeSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","metric.__name__":"up","metric.instance":"localhost:9091","metric.job":"node","status":"success","timestamp":1435781430,"value":0}`,
		`{"baseUrl":"http://localhost:9090/api/v1/","data.resultType":"matrix","event_type":"promTestQueryRangeSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","metric.__name__":"up","metric.instance":"localhost:9091","metric.job":"node","status":"success","timestamp":1435781445,"value":0}`,
		`{"baseUrl":"http://localhost:9090/api/v1/","data.resultType":"matrix","event_type":"promTestQueryRangeSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","metric.__name__":"up","metric.instance":"localhost:9091","metric.job":"node","status":"success","timestamp":1435781460,"value":1}`,
	}
	config := load.Config{
		Name: "Prom Range Test",
		Global: load.Global{
			BaseURL: "http://localhost:9090/api/v1/",
		},
		APIs: []load.API{
			load.API{
				EventType: "promTestQueryRangeSample",
				File:      "../../test/payloads/pqueryRange.json",
			},
		},
	}

	RunConfig(config)
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestE2E_PrometheusPTargetsMetrics(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestPrometheusPTargetsMetrics", "nri-flex")
	expectedSamples := []string{
		`{"baseUrl":"http://localhost:9090/api/v1/","discoveredLabels.__address__":"127.0.0.1:9090","discoveredLabels.__metrics_path__":"/metrics","discoveredLabels.__scheme__":"http","discoveredLabels.job":"prometheus","event_type":"promTestPTargetsSample","health":"up","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","labels.instance":"127.0.0.1:9090","labels.job":"prometheus","lastError":"","lastScrape":"2017-01-17T15:07:44.723715405+01:00","scrapeUrl":"http://127.0.0.1:9090/metrics","status":"success"}`,
		`{"baseUrl":"http://localhost:9090/api/v1/","discoveredLabels.__address__":"127.0.0.1:9100","discoveredLabels.__metrics_path__":"/metrics","discoveredLabels.__scheme__":"http","discoveredLabels.job":"node","event_type":"promTestPTargetsSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","status":"success"}`,
	}
	config := load.Config{
		Name: "Prom Range Test",
		Global: load.Global{
			BaseURL: "http://localhost:9090/api/v1/",
		},
		APIs: []load.API{
			load.API{
				EventType: "promTestPTargetsSample",
				File:      "../../test/payloads/pTargets.json",
			},
		},
	}

	RunConfig(config)
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestE2E_Etcd(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestEtcd", "nri-flex")
	expectedSamples := []string{
		`{"baseUrl":"http://127.0.0.1:2379/v2/","counts.fail":0,"counts.success":745,"event_type":"etcdLeaderSample","follower.id":"6e3bd23ae5f1eae0","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","latency.average":0.017039507382550306,"latency.current":0.000138,"latency.maximum":1.007649,"latency.minimum":0,"latency.standardDeviation":0.05289178277920594,"leader":"924e2e83e93f2560","myCustomAttr":"theValue"}`,
		`{"baseUrl":"http://127.0.0.1:2379/v2/","counts.fail":0,"counts.success":735,"event_type":"etcdLeaderSample","follower.id":"a8266ecf031671f3","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","latency.average":0.012124141496598642,"latency.current":0.000559,"latency.maximum":0.791547,"latency.minimum":0,"latency.standardDeviation":0.04187900156583733,"leader":"924e2e83e93f2560","myCustomAttr":"theValue"}`,
		`{"baseUrl":"http://127.0.0.1:2379/v2/","event_type":"etcdSelfSample","id":"eca0338f4ea31566","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","leaderInfo.leader":"8a69d5f6b7814500","leaderInfo.startTime":"2014-10-24T13:15:51.186620747-07:00","leaderInfo.uptime":"10m59.322358947s","name":"node3","recvAppendRequestCnt":5944,"recvBandwidthRate":570.6254930219969,"recvPkgRate":9.00892789741075,"sendAppendRequestCnt":0,"startTime":"2014-10-24T13:15:50.072007085-07:00","state":"StateFollower"}`,
		`{"baseUrl":"http://127.0.0.1:2379/v2/","compareAndSwapFail":0,"compareAndSwapSuccess":0,"createFail":0,"createSuccess":2,"deleteFail":0,"deleteSuccess":0,"event_type":"etcdStoreSample","expireCount":0,"getsFail":4,"getsSuccess":75,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","setsFail":2,"setsSuccess":4,"updateFail":0,"updateSuccess":0,"watchers":0}`,
	}
	config := load.Config{
		Name: "Prom Range Test",
		Global: load.Global{
			BaseURL: "http://127.0.0.1:2379/v2/",
		},
		APIs: []load.API{
			load.API{
				EventType: "etcdLeaderSample",
				File:      "../../test/payloads/etcdLeader.json",
				SampleKeys: map[string]string{
					"followerSample": "followers>follower.id",
				},
				CustomAttributes: map[string]string{
					"myCustomAttr": "theValue",
				},
			},
			load.API{
				EventType: "etcdSelfSample",
				File:      "../../test/payloads/etcdSelf.json",
			},
			load.API{
				EventType: "etcdStoreSample",
				File:      "../../test/payloads/etcdStore.json",
			},
		},
	}

	RunConfig(config)
	testSamples(expectedSamples, load.Entity.Metrics, t)
}

func TestE2E_Squid(t *testing.T) {
	load.Refresh()
	i, _ := integration.New(load.IntegrationName, load.IntegrationVersion)
	load.Entity, _ = i.Entity("TestSquid", "nri-flex")

	config := load.Config{
		Name: "squidFlex",
		APIs: []load.API{
			{
				Name: "squidMgrUtilization",
				Commands: []load.Command{
					load.Command{
						Run:     "cat ../../test/payloads/squid-mgr.out",
						SplitBy: " = ",
						LineEnd: 88,
					},
				},
				PluckNumbers: true,
				ValueParser: map[string]string{
					"time": "[0-9]+",
				},
			},
		},
	}

	var jsonOut interface{}
	expectedOutput, _ := ioutil.ReadFile("../../test/payloadsExpected/squidMgrTest.json")
	json.Unmarshal(expectedOutput, &jsonOut)
	expectedDatastore := jsonOut.([]interface{})
	RunConfig(config)

	if len(expectedDatastore) != len(load.Entity.Metrics) {
		t.Errorf("expected %d got %d", len(expectedDatastore), len(load.Entity.Metrics))
	}

	for key, expectedVal := range expectedDatastore[0].(map[string]interface{}) {
		if expectedVal != load.Entity.Metrics[0].Metrics[key] {
			t.Errorf("%v expected %v got %v", key, expectedVal, load.Entity.Metrics[0].Metrics[key])
		}
	}
}
