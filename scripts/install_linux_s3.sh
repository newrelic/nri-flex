#!/usr/bin/env bash
echo "### Automated New Relic Flex Install"

if [ "$EUID" -ne 0 ]; then
    echo "### Please run with sudo or root"
    exit
fi

echo "### Pre clean up"
rm -rf /tmp/nri-flex-linux-*

echo "### Checking Latest Version from S3"
curl -o /tmp/NRI-FLEX-LATEST https://newrelic-flex.s3-ap-southeast-2.amazonaws.com/releases/LATEST
VERSION="$(cat /tmp/NRI-FLEX-LATEST)"

echo "### Latest Version from S3: $VERSION"
echo "### Downloading $VERSION from S3"

curl -o "/tmp/nri-flex-linux-$VERSION.tar.gz" "https://newrelic-flex.s3-ap-southeast-2.amazonaws.com/releases/nri-flex-linux-$VERSION.tar.gz"

echo "### Unpacking $VERSION to /tmp/nri-flex-linux-$VERSION/"
tar -xf "/tmp/nri-flex-linux-$VERSION.tar.gz" -C "/tmp/"

echo "### Stopping NR Infrastructure Agent"
if [ -f /etc/systemd/system/newrelic-infra.service ]; then
    echo "SYSTEMD"
    service newrelic-infra stop
fi
if [ -f /etc/init/newrelic-infra.conf ]; then
    echo "INITCTL"
    initctl stop newrelic-infra
fi

echo "### Copying Files"
cp /tmp/nri-flex-linux-$VERSION/nri-flex-config.yml /etc/newrelic-infra/integrations.d/
cp /tmp/nri-flex-linux-$VERSION/nri-flex-definition.yml /var/db/newrelic-infra/custom-integrations/
cp /tmp/nri-flex-linux-$VERSION/nri-flex /var/db/newrelic-infra/custom-integrations/

echo "### Post clean up"
rm -rf /tmp/nri-flex-linux-*

echo "### Creating Directory Structure"
mkdir -p /var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/
mkdir -p /var/db/newrelic-infra/custom-integrations/flexConfigs/

echo "### Starting NR Infrastructure Agent"
if [ -f /etc/systemd/system/newrelic-infra.service ]; then
    service newrelic-infra start
fi
if [ -f /etc/init/newrelic-infra.conf ]; then
    initctl start newrelic-infra
fi