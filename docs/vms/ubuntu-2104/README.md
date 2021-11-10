Ubuntu 2104 VM with simple nri-flex integration that executes script with single numbers as output.

## Requirements

* [Vagrant](https://www.vagrantup.com/)
* NR Account

## Run VM

Get into working directory:
```shell
cd nri-flex/docs/vms/ubuntu-2104
```

Copy agent configuration and add the license_key:
```shell
  cp provision/files/newrelic-infra.yml.dist provision/files/newrelic-infra.yml
  # edit and add your license_key
```

Spawn VM:
```shell
vagrant up
```

Check data in NR1:

NRQL:
```
SELECT * from RandomNumbersSample SINCE 10 minutes ago
```