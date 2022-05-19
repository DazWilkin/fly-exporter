ARG GOLANG_VERSION=1.18
ARG GOLANG_OPTIONS="CGO_ENABLED=0 GOOS=linux GOARCH=amd64"

ARG PROJECT="fly-exporter"

FROM docker.io/golang:${GOLANG_VERSION} as build

ARG PROJECT

WORKDIR /${PROJECT}

COPY main.go ./
COPY collector ./collector

ARG VERSION=""
ARG COMMIT=""

RUN env ${GOLANG_OPTIONS} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/exporter \
    ./main.go

FROM gcr.io/distroless/base-debian11

LABEL org.opencontainers.image.source https://github.com/DazWilkin/fly-exporter

COPY --from=build /go/bin/exporter /

ENTRYPOINT ["/exporter"]