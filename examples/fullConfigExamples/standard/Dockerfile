FROM newrelic/infrastructure

# add netcat
RUN apk add --update netcat-openbsd && rm -rf /var/cache/apk/*

# create some needed default directories for flex
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/flexConfigs/
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/

# if using container discovery configs uncomment this section
ADD flexConfigs /var/db/newrelic-infra/custom-integrations/flexConfigs/

# copy config/definition/binary over
COPY ./nri-flex-config.yml /etc/newrelic-infra/integrations.d/
COPY ./nri-flex-def-nix.yml /var/db/newrelic-infra/custom-integrations/
COPY ./nri-flex /var/db/newrelic-infra/custom-integrations/