//nolint:lll
package fixtures

var FileTests = []struct {
	Name           string
	FileContent    string
	Config         string
	ExpectedStdout string
}{
	{
		Name:        "Use ignore_output",
		FileContent: `{"id": "ec8f4ea31566","leaderInfo": {"leader": "8a67814500","startTime": 1588232295,"uptime": 3600},"name": "node3","sendAppendRequestCnt": 0,"state": "StateFollower"}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    ignore_output: true
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":0,"flex.counter.EventDropCount":0,"flex.time.elapsedMs":45,"flex.time.endMs":1654704471027,"flex.time.startMs":1654704470982}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use sample_filter to skip metrics",
		FileContent: `{"id": "ec8f4ea31566","leaderInfo": {"leader": "8a67814500","startTime": 1588232295,"uptime": 3600},"name": "node3","sendAppendRequestCnt": 0,"state": "StateFollower"}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    sample_filter:
      - name: node3
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":0,"flex.counter.EventDropCount":0,"flex.time.elapsedMs":47,"flex.time.endMs":1654705788545,"flex.time.startMs":1654705788498}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use sample_include_filter to get key",
		FileContent: `{"usageInfo": [{"quantities": 10,"customerId": "abc"},{"quantities": 20,"customerId": "xyz"}]}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    sample_include_filter:
      - customerId: abc
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"customerId":"abc","env":"production","event_type":"usageInfoSample","quantities":10},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.usageInfoSample":1,"flex.time.elapsedMs":8,"flex.time.endMs":1654706477610,"flex.time.startMs":1654706477602}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use sample_exclude_filter to exclude key",
		FileContent: `{"usageInfo": [{"quantities": 10,"customerId": "abc"},{"quantities": 20,"customerId": "xyz"}]}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    sample_include_filter:
      - customerId: abc
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"customerId":"xyz","env":"production","event_type":"usageInfoSample","quantities":20},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.usageInfoSample":1,"flex.time.elapsedMs":49,"flex.time.endMs":1654706661615,"flex.time.startMs":1654706661566}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Provide start_key and rename_key of json",
		FileContent: `{"id": "ec8f4ea31566","leaderInfo": {"leader": "8a67814500","startTime": 1588232295,"uptime": 3600},"name": "node3","sendAppendRequestCnt": 0,"state": "StateFollower"}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    start_key:
      - leaderInfo
    rename_keys:
      startTime: timestamp
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"env":"production","event_type":"readEtcdSelfLeaderInfoSample","leader":"8a67814500","timestamp":1588232295,"uptime":3600},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.readEtcdSelfLeaderInfoSample":1,"flex.time.elapsedMs":8,"flex.time.endMs":1654699825870,"flex.time.startMs":1654699825862}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use keep_keys with jq",
		FileContent: `{"id": "ec8f4ea31566","leaderInfo": {"leader": "8a67814500","startTime": 1588232295,"uptime": 3600},"name": "node3","sendAppendRequestCnt": 0,"state": "StateFollower"}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    keep_keys:
      - leader
    jq: ".leaderInfo"
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"event_type":"readEtcdSelfLeaderInfoSample","leader":"8a67814500"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.readEtcdSelfLeaderInfoSample":1,"flex.time.elapsedMs":9,"flex.time.endMs":1654703090744,"flex.time.startMs":1654703090735}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use lazy_flatten with math and remove_key",
		FileContent: `{"contacts": [{"name": "batman","number": 911},{"name": "robin","number": 112}]}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    lazy_flatten:
      - contacts
    math:
      sum: ${contacts.flat.0.number} + ${contacts.flat.1.number} + 200
    remove_keys:
      - contacts.flat.1.name
    custom_attributes:
      env: production
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"contacts.flat.0.name":"batman","contacts.flat.0.number":911,"contacts.flat.1.number":112,"env":"production","event_type":"readEtcdSelfLeaderInfoSample","sum":1223},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.readEtcdSelfLeaderInfoSample":1,"flex.time.elapsedMs":57,"flex.time.endMs":1654705302291,"flex.time.startMs":1654705302234}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use snake_to_camel and split_objects",
		FileContent: `{"first":{"id":"eca01566","leader_Info":{"up_time":"10m59.322358947s","abc":{"def":123,"hij":234}},"name":"node1"},"second":{"id":"eca04ea31566","leader_Info":{"up_time":"10m59.322358947s","abc":{"def":123,"hij":234}},"name":"node2"}}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    snake_to_camel: true
    split_objects: true
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"event_type":"Sample","id":"eca01566","leaderInfo.abc.def":123,"leaderInfo.abc.hij":234,"leaderInfo.upTime":"10m59.322358947s","name":"node1","split.id":"first"},{"event_type":"Sample","id":"eca04ea31566","leaderInfo.abc.def":123,"leaderInfo.abc.hij":234,"leaderInfo.upTime":"10m59.322358947s","name":"node2","split.id":"second"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":2,"flex.counter.EventDropCount":0,"flex.counter.Sample":2,"flex.time.elapsedMs":45,"flex.time.endMs":1654761239572,"flex.time.startMs":1654761239527}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use to_lower, strip_keys, value_parser and value_transformer",
		FileContent: `{"id":"eca0338f4ea31566","leaderInfo":{"leader":"a8a69d5f6b7814500","startTime":"2014-10-24T13:15:51.186620747-07:00","uptime":"10m59.322358947s","abc":{"def1":"a:123","def2":"a:234"}},"name":"node3"}`,
		Config: `
name: jsonIntegrationTest
apis:
  - name: readEtcdSelfLeaderInfo
    file: FILE_PATH
    to_lower: true
    strip_keys:
      - leaderInfo>leader
      - name
    value_parser:
      def: "[0-9]+"
    value_transformer:
      id: 12345
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"event_type":"readEtcdSelfLeaderInfoSample","id":12345,"leaderinfo.abc.def1":123,"leaderinfo.abc.def2":234,"leaderinfo.starttime":"2014-10-24T13:15:51.186620747-07:00","leaderinfo.uptime":"10m59.322358947s"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.readEtcdSelfLeaderInfoSample":1,"flex.time.elapsedMs":3,"flex.time.endMs":1654769186564,"flex.time.startMs":1654769186561}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use split_array",
		FileContent: `{"status":1,"appstatus":-128,"statusstring":null,"appstatusstring":null,"results":[{"status":-128,"schema":[{"name":"TIMESTAMP","type":6},{"name":"HOST_ID","type":5},{"name":"HOSTNAME","type":9},{"name":"PERCENT_USED","type":6}],"data":[[1582159853733,0,"7605f6bec898",0],[1582159853733,2,"067ea6fc4c22",0],[1582159853733,1,"62a10d3f45e3",0]]}]}`,
		Config: `
name: example
apis:
  - name: voltdb_cpu
    event_type: voltdb
    file: FILE_PATH
    split_array: true
    set_header: [TIMESTAMP, HOST_ID, HOSTNAME, PERCENT_USED]
    start_key:
      - results>data
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"HOSTNAME":"7605f6bec898","HOST_ID":0,"PERCENT_USED":0,"TIMESTAMP":1582159853733,"event_type":"voltdb"},{"HOSTNAME":"067ea6fc4c22","HOST_ID":2,"PERCENT_USED":0,"TIMESTAMP":1582159853733,"event_type":"voltdb"},{"HOSTNAME":"62a10d3f45e3","HOST_ID":1,"PERCENT_USED":0,"TIMESTAMP":1582159853733,"event_type":"voltdb"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":3,"flex.counter.EventDropCount":0,"flex.counter.voltdb":3,"flex.time.elapsedMs":45,"flex.time.endMs":1654776591203,"flex.time.startMs":1654776591158}],"inventory":{},"events":[]}]}`,
	},
	{
		Name:        "Use lookup",
		FileContent: `[{"brand":"honda","car":"civic"},{"brand":"toyota","car":"supra"},{"brand":"mistsubishi","car":"lancer"}]`,
		Config: `
name: LookUps
apis:
  - name: read
    event_type: read
    file: FILE_PATH

  - name: world
    commands:
      - run: echo "${lookup.read:brand}:${lookup.read:car}"
        split_by: ":"
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","data":[{"metrics":[{"brand":"honda","car":"civic","event_type":"read"},{"brand":"toyota","car":"supra","event_type":"read"},{"brand":"mistsubishi","car":"lancer","event_type":"read"},{"event_type":"worldSample","flex.commandTimeMs":3,"honda":"civic"},{"event_type":"worldSample","flex.commandTimeMs":0,"toyota":"supra"},{"event_type":"worldSample","flex.commandTimeMs":0,"mistsubishi":"lancer"},{"event_type":"flexStatusSample","flex.Hostname":"ubuntu-2004-vm","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":6,"flex.counter.EventDropCount":0,"flex.counter.read":3,"flex.counter.worldSample":3,"flex.time.elapsedMs":22,"flex.time.endMs":1672415688172,"flex.time.startMs":1672415688150}],"inventory":{},"events":[]}]}`,
	},
}
