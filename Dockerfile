FROM golang:alpine AS build-env
WORKDIR /usr/local/go/src/github.com/dotmesh-io/terraform-provider-dotscience
COPY . /usr/local/go/src/github.com/dotmesh-io/terraform-provider-dotscience

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl make

ENV GO111MODULE=on

RUN make install-release

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=build-env /usr/local/go/bin/terraform-provider-dotscience /bin/terraform-provider-dotscience

