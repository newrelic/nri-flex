#!/bin/bash
echo "New Relic Integration Installer"

if [ "$EUID" -ne 0 ]
  then echo "Please run with sudo or as root"
  exit
fi

echo "Stopping NR Infrastructure Agent"
if [ -f /etc/systemd/system/newrelic-infra.service ]; then
    echo "SYSTEMD"
    service newrelic-infra stop
fi
if [ -f /etc/init/newrelic-infra.conf ]; then
    echo "INITCTL"
    initctl stop newrelic-infra
fi

echo "Copying Files"
cp ./nri-flex-config.yml /etc/newrelic-infra/integrations.d/
cp ./nri-flex-def-linux.yml /var/db/newrelic-infra/custom-integrations/
cp ./nri-flex /var/db/newrelic-infra/custom-integrations/

echo "Creating Directory Structure"
mkdir -p /var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/
mkdir -p /var/db/newrelic-infra/custom-integrations/flexConfigs/

# if you want JMX support, remember you will need Java 7+ installed
# cp -avr ./nrjmx /var/db/newrelic-infra/custom-integrations/ 

# this will copy all configs, only take what you need
# cp -avr ./flexConfigs /var/db/newrelic-infra/custom-integrations/

# this will copy all container discovery configs, only take what you need
# cp -avr ./flexContainerDiscovery /var/db/newrelic-infra/custom-integrations/


echo "Starting NR Infrastructure Agent"
if [ -f /etc/systemd/system/newrelic-infra.service ]; then
    service newrelic-infra start
fi
if [ -f /etc/init/newrelic-infra.conf ]; then
    initctl start newrelic-infra
fi