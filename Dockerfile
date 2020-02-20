FROM newrelic/infrastructure:1.5.0

# define license key as below, or copy a newrelic-infra.yml over
# refer to here for more info: https://hub.docker.com/r/newrelic/infrastructure/
ENV NRIA_LICENSE_KEY=1234567890abcdefghijklmnopqrstuvwxyz1234

# disable the newrelic infrastructure agent from performing any additional monitoring
# using forwarder mode will only make it responsible for executing integrations
ENV NRIA_IS_FORWARD_ONLY true 


# create some needed default directories for flex
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/flexConfigs/
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/

# if using container discovery configs uncomment this section
# ADD flexConfigs /var/db/newrelic-infra/custom-integrations/flexConfigs/

# copy config/definition/binary over
COPY ./configs/nri-flex-config-linux.yml /etc/newrelic-infra/integrations.d/
COPY ./configs/nri-flex-def-linux.yml /var/db/newrelic-infra/custom-integrations/nri-flex-definition.yml
COPY ./bin/linux/nri-flex /var/db/newrelic-infra/custom-integrations/nri-flex

# # add netcat if needed
# RUN apk add --update netcat-openbsd && rm -rf /var/cache/apk/*