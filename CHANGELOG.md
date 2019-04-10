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