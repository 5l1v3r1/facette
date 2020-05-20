FROM golang:latest AS build

LABEL maintainer="Facette Developers <dev@facette.io>"

RUN apt-get update && \
    apt-get install --no-install-recommends -y librrd-dev npm && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN go get github.com/rakyll/statik

WORKDIR /go/src/facette.io/facette

COPY . /go/src/facette.io/facette

RUN make GO_BUILD_LDFLAGS="-extldflags -static"

FROM scratch

COPY --from=build /go/src/facette.io/facette/bin/facette /facette

ENTRYPOINT [ "/facette" ]

