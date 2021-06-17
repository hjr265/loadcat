// Copyright 2015 The Loadcat Authors. All rights reserved.

package feline

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/hjr265/loadcat/cfg"
	"github.com/hjr265/loadcat/data"
)

type Feline struct {
	sync.Mutex

	base string
}

func New() *Feline {
	return &Feline{
		base: "",
	}
}

func (f *Feline) SetBase(base string) error {
	_, err := os.Stat(base)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(base, 0777)
		if err != nil {
			return err
		}
		err = nil
	}
	if err != nil {
		return err
	}
	f.base = base
	return nil
}

func (f *Feline) Start() error {
	drv := Drivers[cfg.Current.Core.Driver]
	return drv.Start()
}

func (f *Feline) Commit(bal *data.Balancer) error {
	f.Lock()
	defer f.Unlock()

	dir := filepath.Join(f.base, bal.Id.Hex())
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	drv := Drivers[cfg.Current.Core.Driver]
	err = drv.Generate(dir, bal)
	if err != nil {
		return err
	}
	err = drv.Reload()
	if err != nil {
		return err
	}
	return nil
}

var DefaultFeline = New()

func SetBase(dir string) error {
	return DefaultFeline.SetBase(dir)
}

func Start() error {
	return DefaultFeline.Start()
}

func Commit(bal *data.Balancer) error {
	return DefaultFeline.Commit(bal)
}
