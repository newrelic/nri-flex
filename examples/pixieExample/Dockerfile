FROM newrelic/infrastructure:latest

WORKDIR /usr/local/bin

COPY install-pixie-no-auth-docker.sh .

RUN chmod +x ./install-pixie-no-auth-docker.sh

RUN ./install-pixie-no-auth-docker.sh

RUN px version

