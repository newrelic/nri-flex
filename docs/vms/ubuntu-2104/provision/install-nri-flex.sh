#!/usr/bin/env bash

curl -s https://download.newrelic.com/infrastructure_agent/gpg/newrelic-infra.gpg | sudo apt-key add -
printf "deb https://download.newrelic.com/infrastructure_agent/linux/apt hirsute main" | sudo tee -a /etc/apt/sources.list.d/newrelic-infra.list
apt-get update
apt-get install -y newrelic-infra
