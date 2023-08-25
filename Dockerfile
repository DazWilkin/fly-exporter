ARG GOLANG_VERSION=1.21.0

ARG PROJECT="fly-exporter"

ARG COMMIT
ARG VERSION

ARG GOOS="linux"
ARG GOARCH="amd64"

FROM docker.io/golang:${GOLANG_VERSION} as build

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

ARG GOOS
ARG GOARCH

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /bin/exporter \
    ./main.go

FROM gcr.io/distroless/static

LABEL org.opencontainers.image.source https://github.com/DazWilkin/fly-exporter

COPY --from=build /bin/exporter /

ENTRYPOINT ["/exporter"]
