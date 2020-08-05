# gce-image-exporter

A Prometheus exporter that reports on GCE image stats in a project.

## Metrics

```
gce_image_update_time_seconds
gce_image_update_errors_total
gce_image_total
gce_image_files_timestamp
gce_image_files_bytes
```

## Usage

```
docker run -it --rm deviavir/gce-image-exporter:v0.1 --prometheusx.listen-address=:9113 --project=<project> --time=60s
```
