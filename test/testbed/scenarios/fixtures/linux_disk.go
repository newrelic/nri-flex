//nolint:lll
package fixtures

var DiskTests = []struct {
	Name           string
	Config         string
	ExpectedStdout string
}{

	{
		Name: "Linux file list of the /etc/apt directory",
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
              split: horizontal
              set_header: [FileName]
              regex_match: true
              split_by: (\S+)
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"Unknown-SNAPSHOT","data":[{"metrics":[{"FileName":"apt.conf.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"auth.conf.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"keyrings","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"preferences.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"sources.list","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"sources.list.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"FileName":"trusted.gpg.d","event_type":"LinuxFileList","integration_name":"com.newrelic.nri-flex","integration_version":"Unknown-SNAPSHOT"},{"event_type":"flexStatusSample","flex.Hostname":"551daa101010","flex.IntegrationVersion":"Unknown-SNAPSHOT","flex.counter.ConfigsProcessed":1,"flex.counter.EventCount":7,"flex.counter.EventDropCount":0,"flex.counter.LinuxFileList":7,"flex.time.elapsedMs":50,"flex.time.endMs":1654007859822,"flex.time.startMs":1654007859772}],"inventory":{},"events":[]}]}`,
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
}
