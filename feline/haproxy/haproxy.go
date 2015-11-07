// Copyright 2015 The Loadcat Authors. All rights reserved.

package haproxy

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/coreos/go-systemd/dbus"

	"github.com/hjr265/loadcat/cfg"
	"github.com/hjr265/loadcat/data"
	"github.com/hjr265/loadcat/feline"
)

var TplHAProxyConfFrontend = template.Must(template.New("").Parse(`
{{$bal0 := index .Balancers 0}}
frontend loadcat-{{$bal0.Settings.Port}}
	bind *:{{$bal0.Settings.Port}} {{if eq $bal0.Settings.Protocol "https"}}ssl crt {{.Dir}}/server.pem{{end}}
	mode {{$bal0.Settings.Protocol}}

	{{range $bal := .Balancers}}
	acl host_{{$bal.Id.Hex}} hdr(host) -i {{$bal.Settings.Hostname}}
	acl host_{{$bal.Id.Hex}} hdr(host) -i {{$bal.Settings.Hostname}}:{{$bal.Settings.Port}}
	{{end}}

	reqadd X-Forwarded-Proto:\ {{$bal0.Settings.Protocol}}

	{{range $bal := .Balancers}}
	use_backend {{$bal.Id.Hex}} if host_{{$bal.Id.Hex}}
	{{end}}
`))

var TplHAProxyConfBackend = template.Must(template.New("").Parse(`
backend {{.Balancer.Id.Hex}}
	mode {{.Balancer.Settings.Protocol}}

	{{range $srv := .Balancer.Servers}}
	server {{$srv.Id.Hex}} {{$srv.Settings.Address}} check weight {{$srv.Settings.Weight}} {{if eq $srv.Settings.Availability "available"}}{{else if eq $srv.Settings.Availability "backup"}}backup{{else if eq $srv.Settings.Availability "unavailable"}}disabled{{end}}
	{{end}}
`))

type HAProxy struct {
	sync.Mutex

	Systemd *dbus.Conn
}

func (n HAProxy) Generate(dir string, bal *data.Balancer) error {
	n.Lock()
	defer n.Unlock()

	f, err := os.Create(filepath.Join(dir, "haproxy.cfg"))
	if err != nil {
		return err
	}
	err = TplHAProxyConfBackend.Execute(f, struct {
		Dir      string
		Balancer *data.Balancer
	}{
		Dir:      dir,
		Balancer: bal,
	})
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	if bal.Settings.Protocol == "https" {
		buf := bytes.Buffer{}
		buf.Write(bal.Settings.SSLOptions.Certificate)
		buf.WriteByte('\n')
		buf.Write(bal.Settings.SSLOptions.PrivateKey)
		buf.WriteByte('\n')
		err = ioutil.WriteFile(filepath.Join(dir, "server.pem"), buf.Bytes(), 0666)
		if err != nil {
			return err
		}
	}

	f, err = os.Create(filepath.Join(dir, "..", "haproxy.cfg"))
	if err != nil {
		return err
	}
	bals, err := data.ListBalancers()
	if err != nil {
		return err
	}
	err = TplHAProxyConfFrontend.Execute(f, struct {
		Dir       string
		Balancers []data.Balancer
	}{
		Dir:       dir,
		Balancers: bals,
	})
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func (n HAProxy) Reload() error {
	n.Lock()
	defer n.Unlock()

	switch cfg.Current.Haproxy.Mode {
	case "manual":
		return nil

	case "systemd":
		if n.Systemd == nil {
			c, err := dbus.NewSystemdConnection()
			if err != nil {
				return err
			}
			n.Systemd = c
		}

		ch := make(chan string)
		_, err := n.Systemd.ReloadUnit(cfg.Current.Haproxy.Systemd.Service, "replace", ch)
		if err != nil {
			return err
		}
		<-ch

		return nil

	default:
		return errors.New("unknown HAProxy mode")
	}

	panic("unreachable")
}

func init() {
	feline.Register("haproxy", HAProxy{})
}
