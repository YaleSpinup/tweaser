# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.16.15 as builder

ARG version=0.0.0
ARG prerelease
ARG githash=""
ARG buildstamp=""

ADD . .
RUN buffalo build --static -o /bin/app --ldflags "-X github.com/YaleSpinup/tweaser/tweaser.Version=$version -X github.com/YaleSpinup/tweaser/tweaser.VersionPrerelease=$prerelease -X github.com/YaleSpinup/tweaser/tweaser.githash=$githash -X github.com/YaleSpinup/tweaser/tweaser.buildstamp=$buildstamp"

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache python3
RUN ln -s /usr/bin/python3 /usr/bin/python
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache wget gettext && \
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
