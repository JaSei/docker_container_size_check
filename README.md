#Docker container size nagios/icinga check

Check all containers in host.
Default threasholds are 1GB warning and 10GB critical.

`docker_container_size_check -w 100MB -c 2GB`
