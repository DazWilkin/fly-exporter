i# Prometheus Exporter for Fly

[![build-container](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml)

+ `ghcr.io/dazwilkin/fly-exporter:faa1b67cbe8363111d36daf6a1397ef4ec2174ba`

## Container

```bash
TOKEN="[FLY-TOKEN]"
IMAGE="ghcr.io/dazwilkin/fly-exporter:faa1b67cbe8363111d36daf6a1397ef4ec2174ba"

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


IMAGE="ghcr.io/dazwilkin/fly-exporter:${TAG}"

podman build \
--build-arg=GOLANG_OPTIONS="CGO_ENABLED=0 GOOS=linux ${ARCH}" \
--build-arg=COMMIT=$(git rev-parse HEAD) \
--build-arg=VERSION=$(uname --kernel-release) \
--tag={IMAGE} \
--file=./Dockerfile \
.
```

Then:

```bash
POD="exporter"
IMAGE="ghcr.io/dazwilkin/fly-exporter:${TAG}"

podman run \
--detach --tty --rm \
--pod=${POD} \
--name=fly-exporter \
--env=TOKEN=${TOKEN} \
${IMAGE} \
  --endpoint=0.0.0.0:8080


## [Sigstore](https://www.sigstore.dev)

fly-exporter container images are being signed by Sigstore and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/fly-exporter:faa1b67cbe8363111d36daf6a1397ef4ec2174ba
```

> **NOTE** cosign.pub may be downloaded [here](/cosign.pub)

To install `cosign`:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```

