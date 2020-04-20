# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 1.1.2 - 2020-04-16
### Changed
- `inherit_attributes` now allows to collect nested attributes
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
