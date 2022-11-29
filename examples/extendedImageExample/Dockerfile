ARG NRI_KUBERNETES_VERSION=3.4.2

FROM newrelic/nri-kubernetes:${NRI_KUBERNETES_VERSION} as executable

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
ENV AWS_PAGER=""

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  tini \
  unzip \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && rm -rf awscliv2.zip \
  && ./aws/install \
  && which aws

RUN addgroup --gid 2000 nri-agent \
  && adduser --uid 2000 --gid 2000 --gecos "" --disabled-password --disabled-login nri-agent

USER nri-agent
WORKDIR /home/nri-agent

COPY --chown=nri-agent --from=executable /bin/nri-kubernetes /bin/nri-kubernetes

ENTRYPOINT [ "tini", "--", "/bin/nri-kubernetes" ]
