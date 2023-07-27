# Issue

## Container running

```bash
podman container logs fly-exporter
```
Yields

```bash
2022/05/19 22:57:08 main/Collect:
"caller"={"file":"fly.go","line":70}
"level"=0
"msg"="Details"
"app"={"ID":"healthcheck-server","Name":"healthcheck-server","State":"","Status":"running","Deployed":true,"Hostname":"healthcheck-server.fly.dev",...}
```

## Metrics available

```bash
curl localhost:8080/metrics
```
Yields:
```bash
# HELP build_info A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter
# TYPE build_info counter
build_info{git_commit="faa1b67cbe8363111d36daf6a1397ef4ec2174ba",go_version="go1.18.2",os_version="5.15.32-v8+"} 1
# HELP fly_exporter_apps Total Number of Apps
# TYPE fly_exporter_apps counter
fly_exporter_apps{deployed="true",id="ackal-status",name="ackal-status",org_slug="personal",status="running"} 1
fly_exporter_apps{deployed="true",id="healthcheck-server",name="healthcheck-server",org_slug="personal",status="running"} 1
# HELP start_time Exporter start time in Unix epoch seconds
# TYPE start_time gauge
start_time 1.652975685e+09
```

## API

```bash
# `up` metric
curl \
--silent \
http://192.168.1.134:9090/api/v1/query?query=up
```
Yields:
```JSON
{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","instance":"localhost:9090","job":"prometheus-server"},"value":[1652992714.048,"1"]},{"metric":{"__name__":"up","instance":"localhost:9093","job":"alertmanager"},"value":[1652992714.048,"1"]}]}}
```

```bash
# `fly_exporter_apps` metric
curl \
--silent \
http://192.168.1.134:9090/api/v1/query?query=fly_exporter_apps
```
Yields:
```JSON
{"status":"success","data":{"resultType":"vector","result":[]}}
```

Aha! It appears the data is only available intermittently. Is this true for `gcp-exporter` true?

```bash
curl \
--silent \
http://192.168.1.134:9090/api/v1/query?query=fly_exporter_apps
```
Yields:
```JSON
{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"fly_exporter_apps","deployed":"true","id":"ackal-status","instance":"localhost:8080","job":"fly_exporter","name":"ackal-status","org_slug":"personal","status":"running"},"value":[1653002542.421,"1"]},{"metric":{"__name__":"fly_exporter_apps","deployed":"true","id":"healthcheck-server","instance":"localhost:8080","job":"fly_exporter","name":"healthcheck-server","org_slug":"personal","status":"running"},"value":[1653002542.421,"1"]}]}}
```

Yes, it appears to be just a function of the periodicity