//nolint:lll
package fixtures

var CommandTests = []struct {
	Name           string
	Config         string
	ExpectedStdout string
}{

	{
		Name: "Linux file list of the /etc/apt directory with specified shell",
		Config: `
---
integrations:
  - name: nri-flex
    interval: 300s
    config:
      name: LinuxFileList
      apis:
        - event_type: LinuxFileList
          commands:
            - run: ls -1 /etc/apt/
              shell: sh
              split: horizontal
              set_header: [FileName]
              regex_match: true
              split_by: (\S+)
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"FileName":"apt.conf.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"auth.conf.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"keyrings","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"preferences.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"sources.list","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"sources.list.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"trusted.gpg.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"551daa101010","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":7,"flex.counter.EventDropCount":0,"flex.counter.LinuxFileList":7,"flex.time.elapsedMs":50,"flex.time.endMs":1654007859822,"flex.time.startMs":1654007859772}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Linux file list of the /etc/apt directory removing header",
		Config: `
---
integrations:
  - name: nri-flex
    interval: 300s
    config:
      name: LinuxFileList
      apis:
        - event_type: LinuxFileList
          commands:
            - run: ls -l /etc/apt/
              split: horizontal
              set_header: [Permissions,Type,Owner,Group,Size,Month,Day,Hour,FileName]
              row_start: 1
              regex_match: false
              split_by: \s+
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"Day":28,"FileName":"apt.conf.d","Group":"root","Hour":"12:04","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":213,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":8,"FileName":"auth.conf.d","Group":"root","Hour":"10:22","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":6,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":8,"FileName":"keyrings","Group":"root","Hour":"10:22","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":6,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":8,"FileName":"preferences.d","Group":"root","Hour":"10:22","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":6,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":28,"FileName":"sources.list","Group":"root","Hour":"12:04","Month":"Apr","Owner":"root","Permissions":"-rw-r--r--.","Size":2403,"Type":1,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":8,"FileName":"sources.list.d","Group":"root","Hour":"10:22","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":6,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":28,"FileName":"trusted.gpg.d","Group":"root","Hour":"12:04","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":84,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":7,"flex.counter.EventDropCount":0,"flex.counter.LinuxFileList":7,"flex.time.elapsedMs":16,"flex.time.endMs":1654695871048,"flex.time.startMs":1654695871032}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Linux file list of the /etc/apt directory filtering lines",
		Config: `
---
integrations:
  - name: nri-flex
    interval: 300s
    config:
      name: LinuxFileList
      apis:
        - event_type: LinuxFileList
          commands:
            - run: ls -l /etc/apt/
              split: horizontal
              set_header: [Permissions,Type,Owner,Group,Size,Month,Day,Hour,FileName]
              line_start: 1
              line_end: 3
              regex_match: false
              split_by: \s+
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"Day":28,"FileName":"apt.conf.d","Group":"root","Hour":"12:04","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":213,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"Day":8,"FileName":"auth.conf.d","Group":"root","Hour":"10:22","Month":"Apr","Owner":"root","Permissions":"drwxr-xr-x.","Size":6,"Type":2,"event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":2,"flex.counter.EventDropCount":0,"flex.counter.LinuxFileList":2,"flex.time.elapsedMs":21,"flex.time.endMs":1654696245174,"flex.time.startMs":1654696245153}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Linux filesystem usage",
		Config: `
---
integrations:
  - name: nri-flex
    # interval: 30s
    config:
      name: linuxFilesystem
      apis:
        - name: linuxFilesystem
          commands:
            - run: df -PT -B1 -x tmpfs -x xfs -x vxfs -x btrfs -x ext -x ext2 -x ext3 -x ext4 -x hfs
              split: horizontal
              set_header:
                [
                  fs,
                  fsType,
                  capacityBytes,
                  usedBytes,
                  availableBytes,
                  usedPerc,
                  mountedOn,
                ]
              regex_match: true
              split_by: (\S+.\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.*)
          perc_to_decimal: true
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"availableBytes":"Available","capacityBytes":"1-blocks","event_type":"linuxFilesystemSample","fs":"Filesystem","fsType":"Type","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","mountedOn":"Mounted on","usedBytes":"Used","usedPerc":"Capacity"},{"availableBytes":3847839744,"capacityBytes":41921515520,"event_type":"linuxFilesystemSample","fs":"overlay","fsType":"overlay","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","mountedOn":"/","usedBytes":38073675776,"usedPerc":91},{"event_type":"flexStatusSample","flex.Hostname":"0a8c4028be4f","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":2,"flex.counter.EventDropCount":0,"flex.counter.linuxFilesystemSample":2,"flex.time.elapsedMs":57,"flex.time.endMs":1654009494052,"flex.time.startMs":1654009493995}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Linux filesystem usage mock using printf and splitting header",
		Config: `
---
name: linuxFilesystem
apis:
  - name: linuxFilesystem
    commands:
      - run: printf "Filesystem     Type                1K-blocks     Used Available Use Mounted\nfuse-overlayfs fuse.fuse-overlayfs  40938980 37779284   3159696  93 /\n"
        split: horizontal
        header_split_by: \s+
        row_header: 0
        row_start: 1
        regex_match: false
        split_by: \s+
    perc_to_decimal: true
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"1K-blocks":40938980,"Available":3159696,"Filesystem":"fuse-overlayfs","Mounted":"/","Type":"fuse.fuse-overlayfs","Use":93,"Used":37779284,"event_type":"linuxFilesystemSample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.linuxFilesystemSample":1,"flex.time.elapsedMs":49,"flex.time.endMs":1654698466033,"flex.time.startMs":1654698465984}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Echo message with assert pattern",
		Config: `
---
integrations:
  - name: nri-flex
    config:
      name: EchoHi
      apis:
        - name: echoHi
          event_type: echoMessage
          commands:
            - run: "echo hi:bye"
              split_by: ":"
              assert:
                match: hi
                not_match: foo
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"event_type":"echoMessage","flex.commandTimeMs":3,"hi":"bye","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":1,"flex.counter.EventDropCount":0,"flex.counter.echoMessage":1,"flex.time.elapsedMs":40,"flex.time.endMs":1654696788571,"flex.time.startMs":1654696788531}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Print message, store_variable and store_lookups",
		Config: `
---
integrations:
  - name: nri-flex
    config:
      name: jsonIntegrationTest
      apis:
        - name: post
          commands:
           - run: printf '{"id":123,"node":456}\n'
          store_variables:
            Id: id
          store_lookups:
            nodeId: node
        - name: readIDInfo2
          commands:
            - run: printf '{"different_id":${var:Id}}\n'
        - name: readIDInfo3
          commands:
            - run: printf '{"different_node":${lookup:nodeId}}\n'
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"event_type":"postSample","id":123,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT","node":456},{"different_id":123,"event_type":"readIDInfo2Sample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"different_node":456,"event_type":"readIDInfo3Sample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":3,"flex.counter.EventDropCount":0,"flex.counter.postSample":1,"flex.counter.readIDInfo2Sample":1,"flex.counter.readIDInfo3Sample":1,"flex.time.elapsedMs":165,"flex.time.endMs":1654776136702,"flex.time.startMs":1654776136537}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Print message and lookup",
		Config: `
---
integrations:
  - name: nri-flex
    config:
      name: jsonIntegrationTest
      apis:
        - name: post
          commands:
           - run: printf '[{"id":123},{"id":456}]\n'
        - name: readIDInfo2
          commands:
            - run: printf '{"different_id":${lookup.postSample:id}}\n'
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"event_type":"postSample","id":123,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"postSample","id":456,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"different_id":123,"event_type":"readIDInfo2Sample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"different_id":456,"event_type":"readIDInfo2Sample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":4,"flex.counter.EventDropCount":0,"flex.counter.postSample":2,"flex.counter.readIDInfo2Sample":2,"flex.time.elapsedMs":52,"flex.time.endMs":1654775074784,"flex.time.startMs":1654775074732}],"inventory":{},"events":[]}]}`,
	},
	{
		Name: "Print message, lookup and dedupe_lookups",
		Config: `
---
integrations:
  - name: nri-flex
    config:
      name: jsonIntegrationTest
      apis:
        - name: post
          commands:
           - run: printf '[{"id":123},{"id":123}]\n'
        - name: readIDInfo2
          commands:
            - run: printf '{"different_id":${lookup.postSample:id}}\n'
          dedupe_lookups:
            - id
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"event_type":"postSample","id":123,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"postSample","id":123,"integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"different_id":123,"event_type":"readIDInfo2Sample","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"0e0a965295ba","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":3,"flex.counter.EventDropCount":0,"flex.counter.postSample":2,"flex.counter.readIDInfo2Sample":1,"flex.time.elapsedMs":144,"flex.time.endMs":1654775565835,"flex.time.startMs":1654775565691}],"inventory":{},"events":[]}]}`,
	},
}
