// Docker container size nagios/icinga check
//
// Check all containers on host. Default thresholds are 1GB warning and 10GB critical.
//
//	docker_container_size_check -w 100MB -c 2GB

package main

import (
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/olorin/nagiosplugin"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/alecthomas/units"
)

const VERSION = "0.2.0"

var (
	warn = kingpin.Flag("warn", "Warning treshold for image size").Short('w').Default("1GB").Bytes()
	crit = kingpin.Flag("crit", "Critical treshold for image size").Short('c').Default("10GB").Bytes()
)

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.18", nil, defaultHeaders)
	if err != nil {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("connect to docker: %s", err.Error()))
	}

	options := types.ContainerListOptions{Size: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("list docker containers: %s", err.Error()))
	}

	check := nagiosplugin.NewCheck()
	for _, c := range containers {
		warnLevel := level(c, (int64(*warn)), "CHECK_DOCKER_CONTAINER_SIZE_WARN")
		critLevel := level(c, (int64(*crit)), "CHECK_DOCKER_CONTAINER_SIZE_CRIT")

		if  c.SizeRw >= warnLevel {
			check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("%s-%s", c.ID, c.Image))
		} else if c.SizeRw >= critLevel {
			check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("%s-%s", c.ID, c.Image))
		}

		check.AddPerfDatum(fmt.Sprintf("%s-%s", c.ID, c.Image), "b", (float64)(c.SizeRw), (float64)(warnLevel), (float64)(critLevel))
	}

	defer check.Finish()

	check.AddResult(nagiosplugin.OK, "No large container(s)")
}

func level(c types.Container, level int64, overrideLevelKey string) int64 {
	newLevel, ok := c.Labels[overrideLevelKey]
	if ok {
		containerLevel, err := units.ParseBase2Bytes(newLevel)
		if err == nil {
			level = containerLevel
		}
	}

	return level
}
