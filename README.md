i# Prometheus Exporter for Fly

[![build-container](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/fly-exporter/actions/workflows/build.yml)

+ `ghcr.io/dazwilkin/fly-exporter:9cfd1d79aff26695db44be930db4e9930aeece5c`

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

