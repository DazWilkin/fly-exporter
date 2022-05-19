i# Prometheus Exporter for Fly

[![build-container](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml)

+ `ghcr.io/dazwilkin/fly-exporter:9cfd1d79aff26695db44be930db4e9930aeece5c`

## Container

```bash
TOKEN="[FLY-TOKEN]"
IMAGE="ghcr.io/dazwilkin/fly-exporter:9cfd1d79aff26695db44be930db4e9930aeece5c"

podman run \
--interactive --tty --rm \
--env=TOKEN=${TOKEN} \
--publish=8080:8080 \
${IMAGE} \
  --endpoint=0.0.0.0:8080
```

## Raspberry Pi

```bash
if [ "$(getconf LONG_BIT)" -eq 64 ]
then
  # 64-bit Raspian
  ARCH="GOARCH=arm64"
  TAG="arm64"
else
  # 32-bit Raspian
  ARCH="GOARCH=arm GOARM=7"
  TAG="arm32v7"
fi

podman build \
--build-arg=GOLANG_OPTIONS="CGO_ENABLED=0 GOOS=linux ${ARCH}" \
--build-arg=COMMIT=$(git rev-parse HEAD) \
--build-arg=VERSION=$(uname --kernel-release) \
--tag=ghcr.io/dazwilkin/gcp-exporter:${TAG} \
--file=./Dockerfile \
.
```

## [Sigstore](https://www.sigstore.dev)

fly-exporter container images are being signed by Sigstore and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/fly-exporter:9cfd1d79aff26695db44be930db4e9930aeece5c
```

> **NOTE** cosign.pub may be downloaded [here](/cosign.pub)

To install `cosign`:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```

