// Copyright 2015 The Loadcat Authors. All rights reserved.

package nginx

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/coreos/go-systemd/v22/dbus"

	"github.com/hjr265/loadcat/cfg"
	"github.com/hjr265/loadcat/data"
	"github.com/hjr265/loadcat/feline"
)

var TplNginxConf = template.Must(template.New("").Parse(`
upstream {{.Balancer.Id.Hex}} {
	{{if eq .Balancer.Settings.Algorithm "least-connections"}}
		least_conn;
	{{else if eq .Balancer.Settings.Algorithm "source-ip"}}
		ip_hash;
	{{end}}

	{{range $srv := .Balancer.Servers}}
		server  {{$srv.Settings.Address}} weight={{$srv.Settings.Weight}} {{if eq $srv.Settings.Availability "available"}}{{else if eq $srv.Settings.Availability "backup"}}backup{{else if eq $srv.Settings.Availability "unavailable"}}down{{end}};
	{{end}}
}

server {
	{{if eq .Balancer.Settings.Protocol "http"}}
		listen  {{.Balancer.Settings.Port}};
	{{else if eq .Balancer.Settings.Protocol "https"}}
		listen  {{.Balancer.Settings.Port}} ssl;
	{{end}}
	server_name  {{.Balancer.Settings.Hostname}};

	{{if eq .Balancer.Settings.Protocol "https"}}
		ssl                  on;
		ssl_certificate      {{.Dir}}/server.crt;
		ssl_certificate_key  {{.Dir}}/server.key;
	{{end}}

	location / {
		proxy_set_header  Host $host;
		proxy_set_header  X-Real-IP $remote_addr;
		proxy_set_header  X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header  X-Forwarded-Proto $scheme;

		proxy_pass  http://{{.Balancer.Id.Hex}};

		proxy_http_version  1.1;

		proxy_set_header  Upgrade $http_upgrade;
		proxy_set_header  Connection 'upgrade';
	}
}
`))

type Nginx struct {
	sync.Mutex

	Systemd *dbus.Conn
}

func (n Nginx) Generate(dir string, bal *data.Balancer) error {
	n.Lock()
	defer n.Unlock()

	f, err := os.Create(filepath.Join(dir, "nginx.conf"))
	if err != nil {
		return err
	}
	err = TplNginxConf.Execute(f, struct {
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
		err = ioutil.WriteFile(filepath.Join(dir, "server.crt"), bal.Settings.SSLOptions.Certificate, 0666)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(dir, "server.key"), bal.Settings.SSLOptions.PrivateKey, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n Nginx) Reload() error {
	n.Lock()
	defer n.Unlock()

	switch cfg.Current.Nginx.Mode {
	case "systemd":
		if n.Systemd == nil {
			c, err := dbus.NewSystemdConnection()
			if err != nil {
				return err
			}
			n.Systemd = c
		}

		ch := make(chan string)
		_, err := n.Systemd.ReloadUnit(cfg.Current.Nginx.Systemd.Service, "replace", ch)
		if err != nil {
			return err
		}
		<-ch

		return nil

	default:
		return errors.New("unknown Nginx mode")
	}

	panic("unreachable")
}

func init() {
	feline.Register("nginx", Nginx{})
}
