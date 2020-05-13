# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 1.3.0 - 2020-05-13
### Added
- Official support for Kubernetes
- Official support for Windows (Windows Server 2012, Windows Server 2016, Windows Server 2019)
- Official support for the text file API under Linux
- Official support for `store_variables` function

## 1.2.2 - 2020-05-05
### Changed
- (567fa52) Fixing table of contents links in functions.md
- (2190645) fix: rename_keys respects use of regex substitutions - saves a step

## 1.2.1 - 2020-04-30
### Added
- (e7ef014) ci: Adding support for arm and arm64

## 1.2.0 - 2020-04-30
### Added
- Read v4 agent integrations config format to ensure config_file arg is fully supported
- Experimental Save and Filter Functions:
    - save_output: Saves sample output to JSON file.
    - sample_include_match_all_filter: Similar to sample_include_filter except creates samples only with match all key/values supplied in the filter.
### Changed
- Move lookup creation after addAttribute
- Improve path handling
- Improve command behaviour so it returns error sample on failure

## 1.1.2 - 2020-04-16
### Changed
- `inherit_attributes` now allows collecting nested attributes
- Minor improvements

## 1.1.1 - 2020-03-24
### Changed
- New Relic Insights direct sample sender now inherits proxy configuration
### Fixed
- Several bugs fixed

## 1.1.0 - 2020-03-03
### Changed
- Custom attributes can now be transformed, since they are now added before functions are run
- Documentation improvements

## 1.0.0 - 2020-02-25

- First stable release of Flex

## 0.9.7 - 2020-02-19
### Added
- Added `run_async` option for API segment to support async with `store_lookups`
### Changed
- Improved unit testing and logging
- NRJMX tool path is now parameterized
- `remove_keys` is now case-insensitive
- HTTP connection errors are written to the sample
### Removed
- K8s discovery support has been removed. K8s discovery should be handled by [`New Relic K8s discovery`](https://github.com/newrelic/nri-discovery-kubernetes)

> For alpha releases, see the [releases page](https://github.com/newrelic/nri-flex/releases).
