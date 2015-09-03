// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

type Algorithm string

var Algorithms = []Algorithm{
	"round-robin",
	"least-connections",
	"source-ip",
}

func (a Algorithm) Label() string {
	return AlgorithmLabels[a]
}

var AlgorithmLabels = map[Algorithm]string{
	"round-robin":       "Round-robin",
	"least-connections": "Least Connections",
	"source-ip":         "Source IP",
}
