FROM gobuffalo/buffalo:v0.18.14 as builder

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
