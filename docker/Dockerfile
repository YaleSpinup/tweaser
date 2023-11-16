FROM golang:1.21.3-alpine as gobuff

EXPOSE 8080

ENV GOPROXY=https://proxy.golang.org

RUN apk add --no-cache --upgrade apk-tools \
    && apk add --no-cache bash curl openssl git build-base nodejs npm sqlite sqlite-dev mysql-client vim postgresql libpq postgresql-contrib libc6-compat

# Installing linter
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# Installing Yarn
RUN npm i -g --no-progress yarn \
    && yarn config set yarn-offline-mirror /npm-packages-offline-cache \
    && yarn config set yarn-offline-mirror-pruning true

# Installing buffalo binary
RUN go install github.com/gobuffalo/cli/cmd/buffalo@latest
RUN go install github.com/gobuffalo/buffalo-pop/v3@latest

RUN mkdir /src
WORKDIR /src

FROM gobuff as builder

ARG version=0.0.0
ARG prerelease
ARG githash=""
ARG buildstamp=""

RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download && go version
ADD . .
RUN buffalo build --static -o /bin/app --ldflags "-X github.com/YaleSpinup/tweaser/tweaser.Version=$version -X github.com/YaleSpinup/tweaser/tweaser.VersionPrerelease=$prerelease -X github.com/YaleSpinup/tweaser/tweaser.githash=$githash -X github.com/YaleSpinup/tweaser/tweaser.buildstamp=$buildstamp"

FROM python:3-alpine
RUN apk add --no-cache bash ca-certificates && \
    pip3 install awscli --upgrade && \
    aws --version

WORKDIR /bin/

COPY --from=builder /bin/app .
COPY docker/*_config.sh /bin

ENV GO_ENV=production
ENV PORT=8080
ENV ADDR=0.0.0.0

EXPOSE 8080

RUN chown -R nobody:nogroup .
USER nobody

CMD /bin/import_config.sh && /bin/app migrate; /bin/app
