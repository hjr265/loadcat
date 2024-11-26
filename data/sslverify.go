// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

type SSLVerify string

var SSLVerifys = []SSLVerify{
	"off",
	"on",
}

func (p SSLVerify) Label() string {
	return SSLVerifyLabels[p]
}

var SSLVerifyLabels = map[SSLVerify]string{
	"off":  "OFF",
	"on": "ON",
}
