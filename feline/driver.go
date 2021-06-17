// Copyright 2015 The Loadcat Authors. All rights reserved.

package feline

import (
	"github.com/hjr265/loadcat/data"
)

type Driver interface {
	Generate(string, *data.Balancer) error
	Start() error
	Reload() error
}

var Drivers = map[string]Driver{}

func Register(name string, drv Driver) {
	Drivers[name] = drv
}
