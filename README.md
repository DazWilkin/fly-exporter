# Prometheus Exporter for [Fly.io](https://fly.io)

[![build-container](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/DazWilkin/fly-exporter.svg)](https://pkg.go.dev/github.com/DazWilkin/fly-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/DazWilkin/fly-exporter)](https://goreportcard.com/report/github.com/DazWilkin/fly-exporter)

+ `ghcr.io/dazwilkin/fly-exporter:e83420a0142988b412b38d53e53fe857ef875421`

## Container

```bash
TOKEN="[FLY-TOKEN]"
IMAGE="ghcr.io/dazwilkin/fly-exporter:e83420a0142988b412b38d53e53fe857ef875421"

podman run \
--interactive --tty --rm \
--env=TOKEN=${TOKEN} \
--publish=8080:8080 \
${IMAGE} \
  --endpoint=0.0.0.0:8080
```

## Metrics

```bash
curl http://localhost:8080/metrics
```

Yields:

```
# HELP build_info A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter
# TYPE build_info counter
build_info{git_commit="897f2bbe476e834c9a3a0b53784c5d0360bfb5f9",go_version="go1.18.2",os_version="5.15.32-v8+"} 1
# HELP fly_exporter_app_info Info about Applications
# TYPE fly_exporter_app_info counter
fly_exporter_app_info{deployed="true",id="foo",name="foo",org_slug="personal",status="running"} 1
fly_exporter_app_info{deployed="true",id="foo",name="foo",org_slug="personal",status="running"} 1
# HELP fly_exporter_cert_info Info about Certificates
# TYPE fly_exporter_cert_info counter
fly_exporter_cert_info{app_id="foo",app_name="foo",status="Awaiting certificates"} 1
# HELP start_time Exporter start time in Unix epoch seconds
# TYPE start_time gauge
start_time 1.652975685e+09
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
```

## [Sigstore](https://www.sigstore.dev)

fly-exporter container images are being signed by Sigstore and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/fly-exporter:e83420a0142988b412b38d53e53fe857ef875421
```

> **NOTE** cosign.pub may be downloaded [here](/cosign.pub)

To install `cosign`:

```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```

## Similar Exporters

+ [Prometheus Exporter for Azure](https://github.com/DazWilkin/azure-exporter)
+ [Prometheus Exporter for GCP](https://github.com/DazWilkin/gcp-exporter)
+ [Prometheus Exporter for Koyeb](https://github.com/DazWilkin/koyeb-exporter)
+ [Prometheus Exporter for Linode](https://github.com/DazWilkin/linode-exporter)
+ [Prometheus Exporter for Vultr](https://github.com/DazWilkin/vultr-exporter)

<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>
