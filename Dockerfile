FROM newrelic/infrastructure:latest

# define license key as below, or copy a newrelic-infra.yml over
# refer to here for more info: https://hub.docker.com/r/newrelic/infrastructure/
ENV NRIA_LICENSE_KEY=1234567890abcdefghijklmnopqrstuvwxyz1234

# disable the newrelic infrastructure agent from performing any additional monitoring
# using forwarder mode will only make it responsible for executing integrations
ENV NRIA_IS_FORWARD_ONLY true 

# create default configurations dir
RUN mkdir /etc/newrelic-infra/integrations.d/flexConfigs

# edit file below to add your integrations configuration
COPY ./configs/nri-flex-config-linux.yml /etc/newrelic-infra/integrations.d/

# copy binary to default search path
COPY ./bin/linux/nri-flex /var/db/newrelic-infra/newrelic-integrations/bin/

# # add netcat if needed
# RUN apk add --update netcat-openbsd && rm -rf /var/cache/apk/*

# add your configuration files here
# ADD ./examples/flexConfigs/some-configuration-file.yml /etc/newrelic-infra/integrations.d/flexConfigs/

# run Flex in isolation mode (no Infrastructure agent)
ENTRYPOINT ["/var/db/newrelic-infra/newrelic-integrations/bin/nri-flex"]
# use this to pass arguments to flex in isolation mode
CMD [ "-verbose", "-config_dir", "/etc/newrelic-infra/integrations.d/flexConfigs/"]