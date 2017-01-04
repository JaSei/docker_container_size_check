package main

import (
	"fmt"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/olorin/nagiosplugin"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

const VERSION = "0.1.4"

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
		if c.SizeRw >= (int64)(*crit) {
			check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("%s-%s", c.ID, c.Image))
		} else if c.SizeRw >= (int64)(*warn) {
			check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("%s-%s", c.ID, c.Image))
		}

		check.AddPerfDatum(fmt.Sprintf("%s-%s", c.ID, c.Image), "b", (float64)(c.SizeRw), (float64)(*warn), (float64)(*crit))
	}

	defer check.Finish()

	check.AddResult(nagiosplugin.OK, "No large container(s)")
}
