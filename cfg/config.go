// Copyright 2015 The Loadcat Authors. All rights reserved.

package cfg

import (
	"os"

	"github.com/naoina/toml"
)

var Current = struct {
	Core struct {
		Address string
		Dir     string
		Driver  string
	}
	Nginx struct {
		Mode    string
		Systemd struct {
			Service string
		}
	}
}{}

func LoadFile(name string) error {
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		f, err = os.Create(name)
		if err != nil {
			return err
		}
		err = toml.NewEncoder(f).Encode(Current)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	err = toml.NewDecoder(f).Decode(&Current)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Current.Core.Address = ":26590"
	Current.Core.Driver = "nginx"
	Current.Nginx.Mode = "systemd"
	Current.Nginx.Systemd.Service = "nginx.service"
}
