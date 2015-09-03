// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

type Protocol string

var Protocols = []Protocol{
	"http",
	"https",
}

func (p Protocol) Label() string {
	return ProtocolLabels[p]
}

var ProtocolLabels = map[Protocol]string{
	"http":  "HTTP",
	"https": "HTTPS",
}
