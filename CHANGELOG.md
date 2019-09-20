0.8.1-pre
------
- Add support to allow database queries to process async with `db_async: true`

0.8.0-pre
------
- Remove logger and use logrus natively
- Update build scripts and configs
- Improve decoration of docker attributes
- Add additional timing options

0.7.7-pre
------
- Add URL Pagination Support, see examples/flexConfigs/pagination-example.yml

0.7.6-pre
------
- Improve JMX support
- Improve HTTP logging
- Provide custom "mode" option for metric_parser options supported: prefix, suffix, regex, contains

0.7.5-pre
------
- Container discovery - add support for multiple ips
- Redirect logging to stderr

0.7.4-pre
------
- Inherit parent attributes capability (currently works with start_keys)
- Improve/Fix key renamer logic (lost functionality in previous release)

0.7.3-pre
------
- Bug fix key renamer/replacer (thx prodjito)
- Update examples

0.7.2-pre
------
- Improve self container id check

0.7.1-pre
------
- Don't run process discovery for Lambda
- Add container_discovery environment variable
- Update examples

0.7.0-pre
------
- Update environment variable handling
- Update linux s3 install script

0.6.9-pre
------
- No functional changes
- Configs and README updated

0.6.8-pre
------
- Improve container discovery
- Ignore New Relic containers
- Add synchronous config processing capability
- Don't run process discovery if container
- Add one liner S3 install to README
- Detect Kubernetes

0.6.7-pre
------
- Add AWS KMS & Hashicorp Vault secret capability
- Add Fargate Container Discovery capability
- Add uri_path option to JMX options
- Update Status Sample

0.6.6-pre
------
- Refactor cmd entry point
- Improve Lambda support & add ingest feature
- Update Flex Status Sample

0.6.5-pre
------
- Add SAP HANA Database Support
- Update discovery logic

0.6.4-pre
------
- Improve timestamp support
- Add support to recurse through folders in flexConfigs
- Add git syncing support
- Update flexStatusSample
- Improve lookup checking
- Update examples
- Adjust TLS config - only apply RootCAs if tls config enable = true
- Do not log if VERBOSE env is set
- - Makes NR Infra agent log "cannot handle plugin output"

0.6.3-pre
------
- Assign variable stores if nil
- Store lookups earlier

0.6.2-pre
------
- Update Examples
- Update container discovery ipMode defaulting

0.6.1-pre
------
- Allow regex matching feature to work with cache output
- Update examples and configs
- Fix event limiter counter not being unlocked
- Add CA Support to TLSConfig for HTTP requests

0.6.0-pre
------
- Add net dial functionality
- Add lookup file functionality
- Add ability to run Flex as a Lambda
- Add metric api functionality
- Add NR Infra events & inventory support, and granular entity control
- Add split_objects functionality
- Switch to logrus package
- Improve logging
- Refactor various areas

0.5.5-pre
------
- Add regex multi matching and splitting functionality

0.5.4-pre
------
- Move status counter lock to avoid concurrent map write issue

0.5.3-pre
------
- Allow lazy_flatten to work top/parent level

0.5.2-pre
------
- Update Flex Event Counter to use sync and avoid concurrent map writes
- Rewrite Lookup Processor - more robust, supports multiple lookups and can be used in any part of a config not just url

0.5.1-pre
------
- Additional JSON Handling
- Key prefix functionality for any samples

0.5.0-pre
------
- Add value_transformer
- Deprecate replace_keys, but keep backwards compatibility temporarily
- rename_keys now uses regex matching

0.4.9-pre
------
- Fix Prometheus histogram sum & count metrics

0.4.8-pre
------
- Add math functionality
- Add command timeout configurability
- Add timestamp functionality
- Don't add blank command samples
- Nested custom attributes for commands

0.4.7-pre
------
- Allow commands to process cached url data
- Add command cache & line_start option
- Change line_limit to line_end (to align with line_start)

0.4.6-pre
------
- If content-type header is not returned, attempt to validate if the payload is JSON, and continue to process as normal

0.4.5-pre
------
- Add several TLS options, that can be used Global, or under API (with enable flag)
- Default HTTP InsecureSkipVerify: false

0.4.4-pre
------
- Default HTTP InsecureSkipVerify: true
- Move Internal packages to full import paths

0.4.3-pre
------
- Update Prometheus parsing 
- Update tests
- Refactor
- Fix concurrency map write issue
- Add rename_samples functionality

0.4.2-pre
------
- Add variable store functionality

0.4.1-pre
------
- Improve Prometheus metric parsing
- Update SampleFilter functionality

0.4.0-pre
------
- Update database parser

0.3.9-pre
------
- Add Value Parser
- Add Pluck Numbers
- Fix event_type for databases

0.3.8-pre
------
- Update logging, 
- Deprecate "-force_log", use "-verbose" instead

0.3.7-pre
------
- Update & simplify container discovery
- Add regex_matching functionality

0.3.5-pre
------
- Improve container discovery
- Additional logging

0.3.4-pre
------
- Add Insights support

0.3.3-pre
------
- New Relic Init