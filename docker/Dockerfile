# build stage
FROM golang:alpine AS build-env

RUN apk add --no-cache git

WORKDIR /go
ENV GOPATH=/go
ADD . /go/src/git.yale.edu/spinup/ecs-api/
RUN go get github.com/golang/dep/cmd/dep && cd /go/src/git.yale.edu/spinup/ecs-api && dep ensure
RUN go build -o api src/git.yale.edu/spinup/ecs-api/*.go

# final stage
FROM alpine
MAINTAINER Camden Fisher <camden.fisher@yale.edu>

WORKDIR /app
COPY --from=build-env /go/api /app/api
RUN chmod 555 /app/api

RUN apk add --no-cache bash ca-certificates wget python gettext && \
    wget -nv "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" && \
    unzip awscli-bundle.zip && \
    ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws && \
    /usr/local/bin/aws --version

# Install Deco
ARG DECO_VERSION=0.3.1
ARG DECO_OS=linux
ARG DECO_ARCH=amd64
ADD https://github.com/YaleUniversity/deco/releases/download/v${DECO_VERSION}/deco-v${DECO_VERSION}-${DECO_OS}-${DECO_ARCH} /usr/local/bin/deco
RUN chmod 555 /usr/local/bin/deco && deco version

COPY docker/*_config.sh /app
COPY docker/config.deco.json /app/config/config.json
RUN chmod 600 /app/config/config.json

RUN chown -R nobody:nogroup /app
USER nobody

CMD /app/import_config.sh && /app/api -config /app/config/config.json