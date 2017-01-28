#Docker container size nagios/icinga check

Check all containers on host.
Default thresholds are 1GB warning and 10GB critical.

`docker_container_size_check -w 100MB -c 2GB`

If is need override thresholds in specific container, then should use labels:
* CHECK_DOCKER_CONTAINER_SIZE_WARN=BYTES
* CHECK_DOCKER_CONTAINER_SIZE_CRIT=BYTES

for example

```
docker run --label CHECK_DOCKER_CONTAINER_SIZE_WARN=5GB a_lot_of_data_image
```