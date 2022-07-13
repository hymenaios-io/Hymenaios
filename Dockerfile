########
# BASE #
########
ARG GO_VERSION="1.18.3"
ARG DEBIAN_VERSION="bullseye"
FROM golang:${GO_VERSION}-${DEBIAN_VERSION}

COPY . /build/
WORKDIR /build/

ARG BUILD_VERSION="development"
RUN make BUILD_VERSION=${BUILD_VERSION} go-build
RUN chmod 755 /build/argus


#########
# ARGUS #
#########
ARG DEBIAN_VERSION="bullseye"
FROM debian:${DEBIAN_VERSION}-slim
LABEL maintainer="The Argus Authors <developers@release-argus.io>"
RUN \
    apt-get update && \
    apt-get install ca-certificates -y && \
    apt-get clean

COPY entrypoint.sh /entrypoint.sh
COPY --from=0 /build/argus               /usr/bin/argus
COPY --from=0 /build/config.yml.example  /app/config.yml
COPY --from=0 /build/LICENSE             /LICENSE

RUN \
    useradd -u 911 -U -d /app -s /bin/false argus && \
    mkdir -p \
        /app \
        /app/data
WORKDIR /app

EXPOSE     8080
VOLUME     [ "/app/data" ]
ENTRYPOINT [ "/entrypoint.sh" ]
