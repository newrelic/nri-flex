FROM newrelic/infrastructure

# define license key as below, or copy a newrelic-infra.yml over
# refer to here for more info: https://hub.docker.com/r/newrelic/infrastructure/
ENV NRIA_LICENSE_KEY=1234567890abcdefghijklmnopqrstuvwxyz1234

# add netcat
RUN apk add --update netcat-openbsd && rm -rf /var/cache/apk/*

# create some needed default directories for flex
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/flexConfigs/
RUN mkdir -p /var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/

# if using container discovery configs uncomment this section
# ADD flexContainerDiscovery /var/db/newrelic-infra/custom-integrations/flexContainerDiscovery/

# copy config/definition/binary over
COPY ./configs/nri-flex-config.yml /etc/newrelic-infra/integrations.d/
COPY ./configs/nri-flex-def-linux.yml /var/db/newrelic-infra/custom-integrations/nri-flex-definition.yml
COPY ./bin/linux/flex /var/db/newrelic-infra/custom-integrations/nri-flex

# add kubectl
# ENV KUBE_LATEST_VERSION="v1.13.3"
# RUN apk add --update ca-certificates \
#  && apk add --update -t deps curl \
#  && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_LATEST_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
#  && chmod +x /usr/local/bin/kubectl \
#  && apk del --purge deps \
#  && rm /var/cache/apk/*

# install docker client
# ARG DOCKER_CLI_VERSION="18.09.2"
# ENV DOWNLOAD_URL="https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_CLI_VERSION.tgz"
# RUN apk --update add curl \
#     && mkdir -p /tmp/download \
#     && curl -L $DOWNLOAD_URL | tar -xz -C /tmp/download \
#     && mv /tmp/download/docker/docker /usr/local/bin/ \
#     && rm -rf /tmp/download \
#     && apk del curl \
#     && rm -rf /var/cache/apk/*
