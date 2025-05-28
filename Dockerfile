ARG GOLANG_VERSION=1.24.3

ARG PROJECT="fly-exporter"

ARG COMMIT
ARG VERSION

ARG TARGETOS
ARG TARGETARCH

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

ARG PROJECT

WORKDIR /${PROJECT}

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go ./
COPY collector ./collector
COPY terminal ./terminal 

ARG VERSION
ARG COMMIT

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /bin/exporter \
    ./main.go

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.source=https://github.com/DazWilkin/fly-exporter

COPY --from=build /bin/exporter /

ENTRYPOINT ["/exporter"]
