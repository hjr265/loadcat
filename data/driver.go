// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

type Driver string

var Drivers = []Driver{
	"nginx",
}

func (a Driver) Label() string {
	return DriverLabels[a]
}

var DriverLabels = map[Driver]string{
	"nginx": "Nginx",
}
