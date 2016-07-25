package main

import (
	"code.cloudfoundry.org/bytefmt"
	"flag"
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/olorin/nagiosplugin"
	"golang.org/x/net/context"
)

type config struct {
	warn uint64
	crit uint64
}

func parseFlag() config {
	var warn, crit string
	flag.StringVar(&warn, "w", "1GB", "Warning treshold for image size")
	flag.StringVar(&crit, "c", "10GB", "Critical treshold for image size")

	flag.Parse()

	num_warn, err := bytefmt.ToBytes(warn)
	if err != nil {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("convert warn to bytes: %s", err.Error()))
	}

	num_crit, err := bytefmt.ToBytes(crit)
	if err != nil {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("convert crit to bytes: %s", err.Error()))
	}

	return config{warn: num_warn, crit: num_crit}
}

func main() {
	check_conf := parseFlag()

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
		if (uint64)(c.SizeRw) >= check_conf.crit {
			check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("%s-%s", c.ID, c.Image))
		} else if (uint64)(c.SizeRw) >= check_conf.warn {
			check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("%s-%s", c.ID, c.Image))
		}

		check.AddPerfDatum(fmt.Sprintf("%s-%s", c.ID, c.Image), "b", (float64)(c.SizeRw), (float64)(check_conf.warn), (float64)(check_conf.crit))
	}

	defer check.Finish()

	check.AddResult(nagiosplugin.OK, "No large container(s)")
}
