FROM golang:1.16-buster as build-env

ENV GO111MODULE=on
ENV APP_ENV=production

WORKDIR /app

COPY service/common/go.* service/common/
COPY service/command-run/go.* service/command-run/
COPY service/job/go.* service/job/
COPY service/result/go.* service/result/
COPY service/trigger/go.* service/trigger/
COPY Makefile .

RUN make dependencies

# Copy all files
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN make build

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

FROM alpine
COPY --from=build-env /app/cmd /app
RUN apk --no-cache add curl

ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV APP_ENV=production
