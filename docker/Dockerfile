# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.13.6 as builder

ARG version=0.0.0
ARG prerelease
ARG githash=""
ARG buildstamp=""

RUN mkdir -p $GOPATH/src/github.com/YaleSpinup/tweaser
WORKDIR $GOPATH/src/github.com/YaleSpinup/tweaser

ADD . .
RUN dep ensure
RUN buffalo build --static -o /bin/app --ldflags "-X actions.Version=$version -X actions.VersionPrerelease=$prerelease -X actions.githash=$githash -X actions.buildstamp=$buildstamp"

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache wget python gettext && \
    wget -nv "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" && \
    unzip awscli-bundle.zip && \
    ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws && \
    /usr/local/bin/aws --version

WORKDIR /bin/

COPY --from=builder /bin/app .
COPY docker/*_config.sh /bin

ENV GO_ENV=development
ENV PORT=8080
ENV ADDR=0.0.0.0

EXPOSE 8080

CMD /bin/import_config.sh && /bin/app migrate; /bin/app
