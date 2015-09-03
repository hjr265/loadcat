// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

type CipherSuite string

var CipherSuites = []CipherSuite{
	"recommended",
	"legacy",
}

func (c CipherSuite) Label() string {
	return CipherSuiteLabels[c]
}

var CipherSuiteLabels = map[CipherSuite]string{
	"recommended": "Recommended",
	"legacy":      "Legacy",
}
