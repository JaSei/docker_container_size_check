package main

import (
	"github.com/alecthomas/units"
	"github.com/docker/docker/api/types"
	"github.com/olorin/nagiosplugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevel(t *testing.T) {
	container := types.Container{}
	container.Labels = make(map[string]string)
	container.Labels["TEST"] = "10MB"
	container.Labels["INVALID"] = "10XX"

	assert.Equal(t, int64(10), threshold(container, 10, ""), "global threshold setting")
	assert.Equal(t, int64(10*units.Mebibyte), threshold(container, 10, "TEST"), "override threshold setting")
	assert.Equal(t, int64(10), threshold(container, 10, "INVALID"), "invalid override units")
}

func TestCalcLevel(t *testing.T) {
	assert.Equal(t, nagiosplugin.OK, calcLevel(types.Container{SizeRw: 1}, 10, 20))
	assert.Equal(t, nagiosplugin.WARNING, calcLevel(types.Container{SizeRw: 10}, 10, 20))
	assert.Equal(t, nagiosplugin.CRITICAL, calcLevel(types.Container{SizeRw: 20}, 10, 20))
}
